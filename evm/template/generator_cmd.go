package template

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

// commandsModify modifies the application commands.go to use EVM.
func commandsModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd", "commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("cosmosevmcmd", "github.com/cosmos/evm/client"),
			xast.WithNamedImport("cosmosevmserver", "github.com/cosmos/evm/server"),
		)
		if err != nil {
			return err
		}

		content, err = xast.RemoveImports(
			content,
			xast.WithImport("github.com/cosmos/cosmos-sdk/client/keys"),
		)
		if err != nil {
			return err
		}

		// replace server.AddCommandsWithStartCmdOptions with cosmosevmserver.AddCommands
		content, err = xast.ModifyFunction(content, "initRootCmd",
			xast.RemoveFuncCall("server.AddCommandsWithStartCmdOptions"),
			xast.AppendFuncAtLine(`// add Cosmos EVM' flavored TM commands to start server, etc.
		cosmosevmserver.AddCommands(
			rootCmd,
			cosmosevmserver.NewDefaultStartOptions(func(l log.Logger, d dbm.DB, w io.Writer, ao servertypes.AppOptions) cosmosevmserver.Application {
				return newApp(l, d, w, ao).(cosmosevmserver.Application)
			}, app.DefaultNodeHome),
			appExport,
			addModuleInitFlags,
		)`,
				1,
			),
		)
		if err != nil {
			return err
		}

		// add keyring commands
		content = strings.Replace(
			content,
			"keys.Commands()",
			"cosmosevmcmd.KeyCommands(app.DefaultNodeHome, false)",
			1,
		)

		return r.File(genny.NewFileS(cmdPath, content))
	}
}

// rootModify modifies the application root.go to use EVM.
func rootModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd", "root.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("cosmosevmkeyring", "github.com/cosmos/evm/crypto/keyring"),
			xast.WithNamedImport("ibctransferevm", "github.com/cosmos/evm/x/ibc/transfer"),
			xast.WithNamedImport("ibctransfer", "github.com/cosmos/ibc-go/v10/modules/apps/transfer"),
			xast.WithNamedImport("ibctransfertypes", "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"),
		)
		if err != nil {
			return err
		}

		// add module manual registration
		content, err = xast.ModifyFunction(content, "NewRootCmd",
			xast.AppendInsideFuncCall(
				"depinject.Configs", // inject custom msg signer
				"depinject.Provide(app.ProvideMsgEthereumTxCustomGetSigner)",
				-1,
			),
			xast.AppendFuncAtLine(`
				// Since the EVM modules don't support dependency injection, we need to
				// manually register the modules on the client side.
				// This needs to be removed after EVM supports App Wiring.
				evmModules := app.RegisterEVM(clientCtx.Codec, clientCtx.InterfaceRegistry)
				for name, mod := range evmModules {
					moduleBasicManager[name] = module.CoreAppModuleBasicAdaptor(name, mod)
					autoCliOpts.Modules[name] = mod
				}
				// Fix basic manager transfer wrapping
				moduleBasicManager[ibctransfertypes.ModuleName] = ibctransferevm.AppModuleBasic{
					AppModuleBasic: &ibctransfer.AppModuleBasic{},
				}`,
				5),
		)
		if err != nil {
			return err
		}

		// wire client context options
		content = strings.ReplaceAll(content,
			"WithViper(app.Name) // env variable prefix",
			"WithViper(app.Name).WithKeyringOptions(cosmosevmkeyring.Option()).WithLedgerHasProtobuf(true)",
		)

		return r.File(genny.NewFileS(cmdPath, content))
	}
}
