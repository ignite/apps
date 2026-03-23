package template

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

// ibcModify modifies the application ibc.go to use EVM.
func ibcModify(appPath string) genny.RunFn {
	return func(r *genny.Runner) error {
		ibcPath := filepath.Join(appPath, "app", "ibc.go")
		f, err := r.Disk.Find(ibcPath)
		if err != nil {
			return err
		}

		content := f.String()

		content, err = xast.ModifyFunction(content, "registerIBCModules",
			xast.AppendInsideFuncCall("ibctransferkeeper.NewKeeper", "app.Erc20Keeper", 8),
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(ibcPath, content))
	}
}
