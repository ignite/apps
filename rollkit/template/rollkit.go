package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gomodule"
	"github.com/ignite/cli/v28/ignite/pkg/xast"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/templates/field/plushhelpers"
)

// NewRollKitGenerator returns the generator to scaffold a rollkit integration inside an app.
func NewRollKitGenerator(chain *chain.Chain) (*genny.Generator, error) {
	g := genny.New()
	ctx := plush.NewContext()
	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))

	binaryName, err := chain.Binary()
	if err != nil {
		return nil, err
	}

	appPath := chain.AppPath()

	if err := updateDependencies(appPath); err != nil {
		return nil, fmt.Errorf("failed to update go.mod: %w", err)
	}

	g.RunFn(commandsModify(appPath, binaryName, chain.Version))

	return g, nil
}

// commandsModify modifies the application start to use rollkit.
func commandsModify(appPath, binaryName string, version cosmosver.Version) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd/commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		if strings.Contains(f.String(), RollkitV0XStartHandler) {
			return errors.New("rollkit v0.x is already installed. Please remove it before installing rollkit v1.x")
		}

		if strings.Contains(f.String(), RollkitV1XStartHandler) {
			return errors.New("rollkit is already installed.")
		}

		if version.LT(cosmosver.StargateFiftyVersion) {
			return errors.New("rollkit requires Ignite v28+ / Cosmos SDK v0.50+")
		}

		content, err := xast.AppendImports(
			f.String(),
			xast.WithLastNamedImport("abciserver", "github.com/rollkit/go-execution-abci/server"), // TODO(@julienrbrt): Download a specific version via go get beforehand
		)
		if err != nil {
			return err
		}

		// replace potential legacy boilerplate present in an ignite v28 chain.
		content = replaceLegacyAddCommands(content)

		// modify the add commands arguments using xast.
		content, err = xast.ModifyCaller(content, ServerAddCommandsWithStartCmdOptions, func(args []string) ([]string, error) {
			return []string{
				"rootCmd",
				"app.DefaultNodeHome",
				"newApp",
				"appExport",
				`server.StartCmdOptions{
				AddFlags: func(cmd *cobra.Command) {
					abciserver.AddFlags(cmd)
				},
				StartCommandHandler: abciserver.StartHandler(),
			}`,
			}, nil
		})

		return r.File(genny.NewFileS(cmdPath, content))
	}
}

// updateDependencies makes sure the correct dependencies are added to the go.mod files.
// go-execution-abci expects rollkit v1.0 to be used.
func updateDependencies(appPath string) error {
	gomod, err := gomodule.ParseAt(appPath)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod: %w", err)
	}

	gomod.AddNewRequire(GoExecPackage, GoExecVersion, false)
	gomod.AddNewRequire(RollkitPackage, RollkitVersion, true)

	// temporarily add a replace for rollkit
	// it can be removed once we have a tag
	gomod.AddReplace(RollkitPackage, "", RollkitPackage, RollkitVersion)

	// save go.mod
	data, err := gomod.Format()
	if err != nil {
		return fmt.Errorf("failed to format go.mod: %w", err)
	}

	return os.WriteFile(filepath.Join(appPath, "go.mod"), data, 0o644)
}

// replaceLegacyAddCommands replaces the legacy `AddCommands` with a temporary `AddCommandsWithStartCmdOptions` boilerplate.
// Atfterwards, we let the same xast function replace the `AddCommandsWithStartCmdOptions` argument.
func replaceLegacyAddCommands(content string) string {
	return strings.Replace(content, "server.AddCommands(", ServerAddCommandsWithStartCmdOptions+"(", 1)
}
