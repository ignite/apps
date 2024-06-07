package registry

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/google/go-github/v56/github"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"golang.org/x/mod/modfile"
)

const (
	appYMLFileName = "app.ignite.yml"
)

var githubRepoPattern = regexp.MustCompile(`^github\.com\/([a-zA-Z0-9\-_]+)\/([a-zA-Z0-9\-_]+)(\/[a-zA-Z0-9\-_]+)*$`)

// AppRepositoryDetails represents the details of an Ignite app repository.
type AppRepositoryDetails struct {
	PackageURL  string
	Name        string
	Owner       string
	Description string
	Tags        []string
	Stars       int
	License     string
	UpdatedAt   time.Time
	URL         string
	App         AppDetails
}

// AppDetails represents the details of an Ignite app.
type AppDetails struct {
	Name          string
	Description   string
	Path          string
	GoVersion     string
	IgniteVersion string
}

// GetAppDetails returns the details of an Ignite app repository.
func (r Querier) GetAppDetails(ctx context.Context, appName string) (*AppRepositoryDetails, error) {
	apps, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var appEntry AppEntry
	for _, app := range apps {
		if app.Name == appName {
			appEntry = app
		}
	}

	if appEntry.Name == "" {
		return nil, errors.Errorf("app %s not found", appName)
	}

	repoOwner, repoName, err := validatePackageURL(appEntry.RepositoryURL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid package URL")
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
		if name != appName {
			continue
		}

		goMod, err := r.getGoMod(ctx, repo, path.Clean(info.Path))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get go.mod for app %s", name)
		}

		appDetails = AppDetails{
			Name:          name,
			Description:   info.Description,
			Path:          info.Path,
			GoVersion:     goMod.Go.Version,
			IgniteVersion: findCLIVersion(goMod),
		}
	}

	result := &AppRepositoryDetails{
		PackageURL:  appEntry.RepositoryURL,
		Name:        repo.GetName(),
		Owner:       repo.GetOwner().GetLogin(),
		Description: repo.GetDescription(),
		Tags:        repo.Topics,
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

func validatePackageURL(pkgURL string) (owner, name string, err error) {
	parts := githubRepoPattern.FindStringSubmatch(pkgURL)
	if len(parts) != 2 {
		return "", "", errors.Errorf("invalid package URL: %s", pkgURL)
	}

	return parts[1], parts[2], nil
}

func findCLIVersion(modFile *modfile.File) string {
	for _, require := range modFile.Require {
		if strings.HasPrefix(require.Mod.Path, igniteCLIPackage) {
			return require.Mod.Version
		}
	}

	return ""
}
