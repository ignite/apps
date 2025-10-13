package template

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

// appModify modifies the application app.go to use EVM.
func appModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		appGoPath := filepath.Join(appPath, module.PathAppGo)
		f, err := r.Disk.Find(appGoPath)
		if err != nil {
			return err
		}

		// change imports
		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("feegrantkeeper", "cosmossdk.io/x/feegrant/keeper"),
			xast.WithNamedImport("_", "github.com/ethereum/go-ethereum/eth/tracers/js"),
			xast.WithNamedImport("_", "github.com/ethereum/go-ethereum/eth/tracers/native"),
			xast.WithImport("github.com/spf13/cast"),
			xast.WithNamedImport("evmsrvflags", "github.com/cosmos/evm/server/flags"),
			xast.WithNamedImport("erc20keeper", "github.com/cosmos/evm/x/erc20/keeper"),
			xast.WithNamedImport("feemarketkeeper", "github.com/cosmos/evm/x/feemarket/keeper"),
			xast.WithNamedImport("ibctransferkeeper", "github.com/cosmos/evm/x/ibc/transfer/keeper"),
			xast.WithNamedImport("evmkeeper", "github.com/cosmos/evm/x/vm/keeper"),
			xast.WithNamedImport("evmante", "github.com/cosmos/evm/ante"),
			xast.WithNamedImport("evmmempool", "github.com/cosmos/evm/mempool"),
		)
		if err != nil {
			return err
		}

		// remove import
		content, err = xast.RemoveImports(content,
			xast.WithNamedImport("ibctransferkeeper", "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"),
		)
		if err != nil {
			return err
		}

		// change to ethereum coin type
		content = strings.Replace(
			content,
			"ChainCoinType = 60",
			"ChainCoinType = 118",
			1,
		)

		// append modules
		content, err = xast.ModifyStruct(
			content,
			"App",
			xast.AppendStructValue(
				"clientCtx",
				"client.Context",
			),
			xast.AppendStructValue(
				"pendingTxListeners",
				"[]evmante.PendingTxListener",
			),
			xast.AppendStructValue(
				"FeeGrantKeeper",
				"feegrantkeeper.Keeper",
			),
			xast.AppendStructValue(
				"FeeMarketKeeper",
				"feemarketkeeper.Keeper",
			),
			xast.AppendStructValue(
				"EVMKeeper",
				"*evmkeeper.Keeper",
			),
			xast.AppendStructValue(
				"Erc20Keeper",
				"erc20keeper.Keeper",
			),
			xast.AppendStructValue(
				"EVMMempool",
				"*evmmempool.ExperimentalEVMMempool",
			),
		)
		if err != nil {
			return err
		}

		// modify new app function
		content, err = xast.ModifyFunction(
			content,
			"New",
			xast.AppendFuncCodeAtLine(
				`// evm must be instantiated before IBC modules
				if err := app.registerEVMModules(appOpts); err != nil {
					panic(err)
				}`,
				5,
			),
			xast.AppendFuncCodeAtLine(
				`if err := app.postRegisterEVMModules(); err != nil {
					panic(err)
				}`,
				7,
			),
			xast.AppendInsideFuncCall(
				"depinject.Configs", // inject custom msg signer
				"depinject.Provide(ProvideMsgEthereumTxCustomGetSigner)",
				-1,
			),
			xast.AppendInsideFuncCall(
				"depinject.Inject", // inject feegrant keeper via depinject
				"&app.FeeGrantKeeper",
				-1,
			),
			xast.AppendFuncCodeAtLine(
				`// set ante handlers
				maxGasWanted := cast.ToUint64(appOpts.Get(evmsrvflags.EVMMaxTxGasWanted))
				app.setAnteHandler(app.txConfig, maxGasWanted)
				// set evm mempool
				app.setEVMMempool()`,
				12,
			),
		)
		if err != nil {
			return err
		}

		// Add new function
		content, err = xast.AppendFunction(content, `// GetStoreKeysMap returns a map of store keys.
func (app *App) GetStoreKeysMap() map[string]*storetypes.KVStoreKey {
			storeKeysMap := make(map[string]*storetypes.KVStoreKey)
			for _, storeKey := range app.GetStoreKeys() {
				kvStoreKey, ok := app.UnsafeFindStoreKey(storeKey.Name()).(*storetypes.KVStoreKey)
				if ok {
					storeKeysMap[storeKey.Name()] = kvStoreKey
				}
			}

			return storeKeysMap
}`)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(appGoPath, content))
	}
}
