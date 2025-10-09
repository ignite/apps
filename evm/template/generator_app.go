package template

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

// appModify modifies the application app.go to use EVM.
func appModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd/commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("_", "github.com/ethereum/go-ethereum/eth/tracers/js"),
			xast.WithNamedImport("_", "github.com/ethereum/go-ethereum/eth/tracers/native"),
			xast.WithImport("github.com/spf13/cast"),
			xast.WithNamedImport("evmsrvflags", "github.com/cosmos/evm/server/flags"),
			xast.WithNamedImport("erc20keeper", "github.com/cosmos/evm/x/erc20/keeper"),
			xast.WithNamedImport("feemarketkeeper", "github.com/cosmos/evm/x/feemarket/keeper"),
			xast.WithNamedImport("ibctransferkeeper", "github.com/cosmos/evm/x/ibc/transfer/keeper"),
			xast.WithNamedImport("evmkeeper", "github.com/cosmos/evm/x/vm/keeper"),
		)
		if err != nil {
			return err
		}

		// add keyring commands
		content = strings.Replace(
			content,
			"keys.Commands()",
			"cosmosevmcmd.KeyCommands(app.DefaultNodeHome, true)",
			1,
		)

		return r.File(genny.NewFileS(cmdPath, content))
	}
}
