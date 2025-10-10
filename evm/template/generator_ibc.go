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

		// remove imports that are replaced by EVM wrappers
		content, err := xast.RemoveImports(f.String(),
			xast.WithNamedImport("ibctransferkeeper", "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"),
			xast.WithNamedImport("ibctransferv2", "github.com/cosmos/ibc-go/v10/modules/apps/transfer/v2"),
		)
		if err != nil {
			return err
		}

		// change imports
		content, err = xast.AppendImports(content,
			xast.WithNamedImport("ibctransferevm", "github.com/cosmos/evm/x/ibc/transfer"),
			xast.WithNamedImport("ibctransferkeeper", "github.com/cosmos/evm/x/ibc/transfer/keeper"),
			xast.WithNamedImport("ibctransferv2evm", "github.com/cosmos/evm/x/ibc/transfer/v2"),
			xast.WithNamedImport("erc20", "github.com/cosmos/evm/x/erc20"),
			xast.WithNamedImport("erc20v2", "github.com/cosmos/evm/x/erc20/v2"),
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(ibcPath, content))
	}
}
