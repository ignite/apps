package registry

import (
	"context"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/google/go-github/v56/github"
	"github.com/iancoleman/strcase"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"golang.org/x/mod/modfile"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// ValidateApp validates the app yaml file.
func (r Querier) ValidateApp(ctx context.Context, appName string) (*AppRepositoryDetails, error) {
	apps, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var appEntry App
	for _, app := range apps {
		if strings.EqualFold(app.Name, appName) {
			appEntry = app
		}
	}

	if appEntry.Name == "" && appEntry.Slug == "" {
		return nil, errors.Errorf("app %s not found", appName)
	}

	if appEntry.Name == "" {
		return nil, errors.Errorf("app name must be defined")
	} else if appEntry.Name != strings.ToUpper(appEntry.Name) {
		return nil, errors.Errorf("app name must be uppercase: %s", appEntry.Name)
	}

	if appEntry.Slug == "" {
		return nil, errors.Errorf("app slug must be defined")
	} else if appEntry.Slug != strcase.ToKebab(appEntry.Name) && appEntry.Slug != strcase.ToSnake(appEntry.Name) {
		return nil, errors.Errorf("app slug must be kebab or snake case: %s", appEntry.Slug)
	}

	if appEntry.AppDescription == "" {
		return nil, errors.New("app description must be defined")
	}

	if err := validateVersion(appEntry.Ignite); err != nil {
		return nil, errors.Wrap(err, "invalid ignite version")
	}

	if err := validateVersion(appEntry.CosmosSDK); err != nil {
		return nil, errors.Wrap(err, "invalid cosmos-sdk version")
	}

	for depName, depVersion := range appEntry.Dependencies {
		if err := validateVersion(depVersion); err != nil {
			return nil, errors.Wrapf(err, "invalid %s version", depName)
		}
	}

	for _, author := range appEntry.Authors {
		if err := validateAuthor(author); err != nil {
			return nil, errors.Wrap(err, "invalid author")
		}
	}

	if err := validateURL(appEntry.RepositoryURL, "repository URL"); err != nil {
		return nil, err
	}

	if err := validateURL(appEntry.DocumentationURL, "documentation URL"); err != nil {
		return nil, err
	}

	if err := validateURL(appEntry.License.URL, "license URL"); err != nil {
		return nil, err
	}

	if appEntry.License.Name == "" {
		return nil, errors.New("license name must be defined")
	}

	if len(appEntry.Keywords) == 0 {
		return nil, errors.New("keywords must be defined")
	}

	if err := validatePlatforms(appEntry.SupportedPlatforms); err != nil {
		return nil, err
	}

	repoOwner, repoName, err := validateRepoURL(appEntry.RepositoryURL)
	if err != nil {
		return nil, err
	}

	repo, err := r.client.GetRepository(ctx, repoOwner, repoName)
	if err != nil {
		return nil, err
	}

	appYML, err := r.getAppsConfig(ctx, repo)
	if err != nil {
		return nil, err
	}

	var appDetails AppDetails
	for name, info := range appYML.Apps {
		if !strings.EqualFold(name, appName) {
			continue
		}

		goMod, err := r.getGoMod(ctx, repo, path.Clean(info.Path))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get go.mod for app %s", name)
		}

		appDetails = AppDetails{
			Name:             name,
			PackageURL:       path.Join(stripHTTPOrHTTPSFromURL(appEntry.RepositoryURL), info.Path),
			DocumentationURL: appEntry.DocumentationURL,
			Description:      info.Description,
			Path:             info.Path,
			GoVersion:        goMod.Go.Version,
			IgniteVersion:    findCLIVersion(goMod),
		}
	}

	result := &AppRepositoryDetails{
		Name:        repo.GetName(),
		Description: repo.GetDescription(),
		Stars:       repo.GetStargazersCount(),
		License:     repo.GetLicense().GetName(),
		UpdatedAt:   repo.GetUpdatedAt().Time,
		URL:         repo.GetHTMLURL(),
		App:         appDetails,
	}

	return result, nil
}

func (r Querier) getGoMod(ctx context.Context, repo *github.Repository, fpath string) (*modfile.File, error) {
	contents, err := r.client.GetFileContent(ctx, repo.GetOwner().GetLogin(), repo.GetName(), path.Join(fpath, "go.mod"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file content")
	}

	mod, err := modfile.Parse(fmt.Sprintf("%s/%s", repo.GetFullName(), "go.mod"), contents, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse go.mod")
	}

	return mod, nil
}

func (r Querier) getAppsConfig(ctx context.Context, repo *github.Repository) (*plugin.AppsConfig, error) {
	data, err := r.client.GetFileContent(ctx, repo.GetOwner().GetLogin(), repo.GetName(), appYMLFileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s file content", appYMLFileName)
	}

	var conf plugin.AppsConfig
	if err := yaml.UnmarshalContext(ctx, data, &conf); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal %s file", appYMLFileName)
	}
	return &conf, nil
}

func validateRepoURL(repoURL string) (owner, name string, err error) {
	matches := githubRepoPattern.FindStringSubmatch(repoURL)
	if len(matches) < 4 {
		return "", "", errors.Errorf("invalid repo URL: %s", repoURL)
	}

	return matches[2], matches[3], nil
}

func findCLIVersion(modFile *modfile.File) string {
	for _, require := range modFile.Require {
		if strings.HasPrefix(require.Mod.Path, igniteCLIPackage) {
			return require.Mod.Version
		}
	}

	return ""
}

// stripHTTPOrHTTPSFromURL strips http or https scheme from a URL.
func stripHTTPOrHTTPSFromURL(url string) string {
	if url[:8] == "https://" {
		url = url[8:]
	} else if url[:7] == "http://" {
		url = url[7:]
	}
	return url
}
