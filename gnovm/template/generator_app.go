package template

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

// appModify modifies the application app.go to use GnoVM.
func appModify(appPath string) genny.RunFn {
	return func(r *genny.Runner) error {
		appGoPath := filepath.Join(appPath, module.PathAppGo)
		f, err := r.Disk.Find(appGoPath)
		if err != nil {
			return err
		}

		// change imports
		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("gnovmmodulekeeper", "github.com/ignite/gnovm/x/gnovm/keeper"),
			xast.WithImport("github.com/cosmos/cosmos-sdk/x/auth/ante"),
		)
		if err != nil {
			return err
		}

		// append modules
		content, err = xast.ModifyStruct(
			content,
			"App",
			xast.AppendStructValue(
				"GnoVMKeeper",
				"gnovmmodulekeeper.Keeper",
			),
		)
		if err != nil {
			return err
		}

		// modify the new app function
		content, err = xast.ModifyFunction(
			content,
			"New",
			xast.AppendInsideFuncCall(
				"depinject.Inject", // inject gnovm keeper via depinject
				"&app.GnoVMKeeper",
				-1,
			),
			xast.AppendFuncAtLine(
				`// set ante handlers
				if err := app.setAnteHandler(ante.HandlerOptions{
					AccountKeeper:   app.AuthKeeper,
					BankKeeper:      app.BankKeeper,
					SignModeHandler: app.txConfig.SignModeHandler(),
					SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
				}); err != nil {
					panic(err)
				}`,
				8,
			),
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(appGoPath, content))
	}
}
