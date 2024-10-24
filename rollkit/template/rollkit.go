package template

import (
	"fmt"
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

// commandsModify modifies the application start to use rollkit.
func commandsModify(appPath, binaryName string, version cosmosver.Version) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd/commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		if strings.Contains(f.String(), "rollserv.StartHandler[servertypes.Application]") {
			return errors.New("rollkit is already installed")
		}

		if version.LT(cosmosver.StargateFiftyVersion) {
			return errors.New("rollkit requires Ignite v28+ / Cosmos SDK v0.50+")
		}
		DefaultRollkitServerVersion := "test-rollkit-main"
		DefaultRollkitConfigVersion := "v0.13.7"
		content, err := xast.AppendImports(
			f.String(),
			xast.WithLastNamedImport("rollserv", fmt.Sprintf("github.com/rollkit/cosmos-sdk-starter/server@%s", DefaultRollkitServerVersion)),
			xast.WithLastNamedImport("rollconf", fmt.Sprintf("github.com/rollkit/rollkit/config@%s", DefaultRollkitConfigVersion)),
		)
		if err != nil {
			return err
		}

		// TODO(@julienrbrt) eventually use ast for simply replacing AddCommands or AddCommandsWithStartCmdOptions
		const (
			defaultv050ServerOptions       = "server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)"
			secondv050DefaultServerOptions = `server.AddCommandsWithStartCmdOptions(rootCmd, app.DefaultNodeHome, newApp, appExport, server.StartCmdOptions{
				AddFlags: func(startCmd *cobra.Command) {
					addModuleInitFlags(startCmd)
				},
			})`
			thirdv050DefaultServerOptions = "server.AddCommandsWithStartCmdOptions(rootCmd, app.DefaultNodeHome, newApp, appExport, server.StartCmdOptions[servertypes.Application]{})"

			rollkitServerOptions = `server.AddCommandsWithStartCmdOptions(
				rootCmd,
				app.DefaultNodeHome,
				newApp, appExport,
				server.StartCmdOptions{
					AddFlags: func(cmd *cobra.Command) {
						rollconf.AddFlags(cmd)
						addModuleInitFlags(cmd)
					},
					StartCommandHandler: rollserv.StartHandler[servertypes.Application],
				})`
		)

		// try all 3 possible default server options
		content = strings.ReplaceAll(content, defaultv050ServerOptions, rollkitServerOptions)
		content = strings.ReplaceAll(content, secondv050DefaultServerOptions, rollkitServerOptions)
		content = strings.ReplaceAll(content, thirdv050DefaultServerOptions, rollkitServerOptions)

		return r.File(genny.NewFileS(cmdPath, content))
	}
}
