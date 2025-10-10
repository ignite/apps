package template

import (
	"path/filepath"
	"strings"

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

		// add EVM imports
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

		// Modify registerIBCModules function
		content, err = modifyRegisterIBCModules(content)
		if err != nil {
			return err
		}

		// Modify RegisterIBC function
		content, err = xast.ModifyFunction(content, "RegisterIBC",
			xast.ReplaceFuncBody(`modules := map[string]appmodule.AppModule{
				ibcexported.ModuleName: ibc.NewAppModule(&ibckeeper.Keeper{}),
				icatypes.ModuleName:    icamodule.NewAppModule(&icacontrollerkeeper.Keeper{}, &icahostkeeper.Keeper{}),
				ibctm.ModuleName:       ibctm.NewAppModule(ibctm.NewLightClientModule(cdc, ibcclienttypes.StoreProvider{})),
				solomachine.ModuleName: solomachine.NewAppModule(solomachine.NewLightClientModule(cdc, ibcclienttypes.StoreProvider{})),
			}

			for _, m := range modules {
				if mr, ok := m.(module.AppModuleBasic); ok {
					mr.RegisterInterfaces(cdc.InterfaceRegistry())
				}
			}

			// manually register for ibctransfer, as instantiation requires a keeper
			ibcTransferModuleBasic := ibctransferevm.AppModuleBasic{
				AppModuleBasic: &ibctransfer.AppModuleBasic{},
			}
			ibcTransferModuleBasic.RegisterInterfaces(cdc.InterfaceRegistry())

			return modules`),
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(ibcPath, content))
	}
}

// modifyRegisterIBCModules modifies the registerIBCModules function to add EVM integration.
func modifyRegisterIBCModules(content string) (string, error) {
	// add app.Erc20Keeper parameter to TransferKeeper initialization
	content, err := xast.ModifyFunction(content, "registerIBCModules",
		xast.AppendInsideFuncCall("ibctransferkeeper.NewKeeper", "app.Erc20Keeper", 8),
	)
	if err != nil {
		return "", err
	}

	// Replace ibctransfer.NewIBCModule with ibctransferevm.NewIBCModule in var declarations
	// Using string replacement for variable initialization as xast doesn't have direct support
	content = strings.ReplaceAll(
		content,
		"ibctransfer.NewIBCModule(app.TransferKeeper)",
		"ibctransferevm.NewIBCModule(app.TransferKeeper)",
	)

	// Replace ibctransferv2.NewIBCModule with ibctransferv2evm.NewIBCModule
	content = strings.ReplaceAll(
		content,
		"ibctransferv2.NewIBCModule(app.TransferKeeper)",
		"ibctransferv2evm.NewIBCModule(app.TransferKeeper)",
	)

	// Add ERC20 middleware wrapping after the var declaration block
	// Find the line after the var block and insert the middleware code
	varBlockEnd := `icaHostStack       porttypes.IBCModule = icahost.NewIBCModule(app.ICAHostKeeper)
	)`

	middlewareCode := `icaHostStack       porttypes.IBCModule = icahost.NewIBCModule(app.ICAHostKeeper)
	)

	// add evm capabilities
	transferStack = erc20.NewIBCMiddleware(app.Erc20Keeper, transferStack)
	transferStackV2 = erc20v2.NewIBCMiddleware(transferStackV2, app.Erc20Keeper)`

	content = strings.Replace(content, varBlockEnd, middlewareCode, 1)

	// Replace ibctransfer.NewAppModule with ibctransferevm.NewAppModule
	content = strings.ReplaceAll(
		content,
		"ibctransfer.NewAppModule(app.TransferKeeper),",
		"ibctransferevm.NewAppModule(app.TransferKeeper), // ibc transfer evm compatible",
	)

	return content, nil
}
