package template

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
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

	g.RunFn(commandsModify(chain.AppPath(), binaryName, chain.Version))

	return g, nil
}

const (
	ServerAddCommandsWithStartCmdOptions = "server.AddCommandsWithStartCmdOptions"

	RollkitV0XStartHandler = "rollserv.StartHandler"
	RollkitV1XStartHandler = "abciserver.StartHandler"
)

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
func updateDependencies() error {
	return nil
}

// replaceLegacyAddCommands replaces the legacy `AddCommands` with a temporary `AddCommandsWithStartCmdOptions` boilerplate.
// Atfterwards, we let the same xast function replace the `AddCommandsWithStartCmdOptions` argument.
func replaceLegacyAddCommands(content string) string {
	return strings.Replace(content, "server.AddCommands(", ServerAddCommandsWithStartCmdOptions+"(", 1)
}
