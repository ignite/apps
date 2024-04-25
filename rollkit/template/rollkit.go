package template

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

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

	g.RunFn(commandsModify(chain.AppPath(), binaryName))

	return g, nil
}

// commandsModify modifies the application start to use rollkit
func commandsModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd/commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(
			f.String(),
			"rollserv github.com/rollkit/cosmos-sdk-starter/server",
			"rollconf github.com/rollkit/rollkit/config",
		)
		if err != nil {
			return err
		}

		const (
			defaultServerOptions = "server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)"
			rollkitServerOptions = `server.AddCommandsWithStartCmdOptions(
				rootCmd,
				app.DefaultNodeHome,
				newApp, appExport,
				server.StartCmdOptions{
					AddFlags:            rollconf.AddFlags,
					StartCommandHandler: rollserv.StartHandler[servertypes.Application],
				})`
		)
		content = strings.ReplaceAll(content, defaultServerOptions, rollkitServerOptions)

		return r.File(genny.NewFileS(cmdPath, content))
	}
}
