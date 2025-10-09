package template

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

//go:embed files/* files/**/*
var fsAppEvm embed.FS

// NewEVMGenerator returns the generator to scaffold a evm integration inside an app.
func NewEVMGenerator(chain *chain.Chain) (*genny.Generator, error) {
	appEvm, err := fs.Sub(fsAppEvm, "files")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()

	if err := g.OnlyFS(appEvm, nil, nil); err != nil {
		return g, err
	}

	appPath := chain.AppPath()
	modpath, _, err := gomodulepath.Find(appPath)
	if err != nil {
		return nil, err
	}

	ctx := plush.NewContext()
	ctx.Set("ModulePath", modpath)
	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))

	binaryName, err := chain.Binary()
	if err != nil {
		return nil, err
	}

	if err := updateDependencies(appPath); err != nil {
		return nil, errors.Errorf("failed to update go.mod: %w", err)
	}

	g.RunFn(commandsModify(appPath, binaryName))
	g.RunFn(rootModify(appPath, binaryName))
	g.RunFn(appModify(appPath, binaryName))

	/// --------------

	// Update existing files
	if err := updateExistingFiles(appPath, binaryName); err != nil {
		return nil, errors.Wrap(err, "failed to update existing files")
	}

	return g, nil
}

// updateDependencies makes sure the correct dependencies are added to the go.mod files.
func updateDependencies(appPath string) error {
	gomod, err := gomodule.ParseAt(appPath)
	if err != nil {
		return errors.Errorf("failed to parse go.mod: %w", err)
	}

	gomod.AddNewRequire(CosmosEVMPackage, CosmosEVMVersion, false)

	// add required replaces
	gomod.AddReplace(EthereumGoEthereumPackage, "", CosmosGoEthereumPackage, CosmosGoEthereumVersion)

	// save go.mod
	data, err := gomod.Format()
	if err != nil {
		return errors.Errorf("failed to format go.mod: %w", err)
	}

	return os.WriteFile(filepath.Join(appPath, "go.mod"), data, 0o644)
}

func updateExistingFiles(appPath, appName string) error {
	// Update app.go
	if err := updateAppGo(appPath, appName); err != nil {
		return err
	}

	// Update app_config.go
	if err := updateAppConfigGo(appPath); err != nil {
		return err
	}

	return nil
}

func updateAppGo(appPath, appName string) error {
	appGoPath := filepath.Join(appPath, "app", "app.go")

	content, err := os.ReadFile(appGoPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string

	importSectionStarted := false
	importSectionEnded := false

	for i, line := range lines {
		// Add EVM imports
		if strings.Contains(line, "import (") && !importSectionStarted {
			importSectionStarted = true
			updatedLines = append(updatedLines, line)
			continue
		}

		if importSectionStarted && !importSectionEnded && strings.TrimSpace(line) == ")" {
			// Add EVM imports before closing import
			updatedLines = append(updatedLines, "")
			updatedLines = append(updatedLines, "\t// Force-load the tracer engines to trigger registration due to Go-Ethereum v1.10.15 changes")
			updatedLines = append(updatedLines, "\t\"github.com/spf13/cast\"")
			updatedLines = append(updatedLines, "")
			updatedLines = append(updatedLines, "\tevmsrvflags \"github.com/cosmos/evm/server/flags\"")
			updatedLines = append(updatedLines, "\terc20keeper \"github.com/cosmos/evm/x/erc20/keeper\"")
			updatedLines = append(updatedLines, "\tfeemarketkeeper \"github.com/cosmos/evm/x/feemarket/keeper\"")
			updatedLines = append(updatedLines, "\tevmkeeper \"github.com/cosmos/evm/x/vm/keeper\"")
			updatedLines = append(updatedLines, "\tibctransferkeeper \"github.com/cosmos/evm/x/ibc/transfer/keeper\"")
			importSectionEnded = true
		}

		// Update ChainCoinType
		line = strings.Replace(line, "ChainCoinType = 118", "ChainCoinType = 60", 1)

		// Add EVM keepers to App struct
		if strings.Contains(line, "TransferKeeper      ibctransferkeeper.Keeper") {
			updatedLines = append(updatedLines, line)
			updatedLines = append(updatedLines, "")
			updatedLines = append(updatedLines, "\t// cosmos evm keepers")
			updatedLines = append(updatedLines, "\tFeeMarketKeeper feemarketkeeper.Keeper")
			updatedLines = append(updatedLines, "\tEVMKeeper       *evmkeeper.Keeper")
			updatedLines = append(updatedLines, "\tErc20Keeper     erc20keeper.Keeper")
			continue
		}

		// Add EVM module registration in New function
		if strings.Contains(line, "if err := app.registerIBCModules(appOpts); err != nil {") {
			updatedLines = append(updatedLines, "\t// evm must be instantiated before IBC modules")
			updatedLines = append(updatedLines, "\tif err := app.registerEVMModules(appOpts); err != nil {")
			updatedLines = append(updatedLines, "\t\tpanic(err)")
			updatedLines = append(updatedLines, "\t}")
			updatedLines = append(updatedLines, "")
			updatedLines = append(updatedLines, line)
			continue
		}

		if strings.Contains(line, "panic(err)") && strings.Contains(lines[i-2], "registerIBCModules") {
			updatedLines = append(updatedLines, line)
			updatedLines = append(updatedLines, "")
			updatedLines = append(updatedLines, "\tif err := app.postRegisterEVMModules(); err != nil {")
			updatedLines = append(updatedLines, "\t\tpanic(err)")
			updatedLines = append(updatedLines, "\t}")
			continue
		}

		// Add ante handler setup
		if strings.Contains(line, "return app.App.InitChainer(ctx, req)") {
			updatedLines = append(updatedLines, line)
			updatedLines = append(updatedLines, "\t})")
			updatedLines = append(updatedLines, "")
			updatedLines = append(updatedLines, "\t// set ante handlers")
			updatedLines = append(updatedLines, "\tmaxGasWanted := cast.ToUint64(appOpts.Get(evmsrvflags.EVMMaxTxGasWanted))")
			updatedLines = append(updatedLines, "\tapp.setAnteHandler(app.txConfig, maxGasWanted)")
			continue
		}

		updatedLines = append(updatedLines, line)

		// quick fix to remove
		strings.ReplaceAll(line, "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper", "")
	}

	return os.WriteFile(appGoPath, []byte(strings.Join(updatedLines, "\n")), 0644)
}

func updateAppConfigGo(appPath string) error {
	configPath := filepath.Join(appPath, "app", "app_config.go")

	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string

	for i, line := range lines {
		// Add EVM module imports
		if strings.Contains(line, "\"google.golang.org/protobuf/types/known/durationpb\"") {
			updatedLines = append(updatedLines, line)
			updatedLines = append(updatedLines, "")
			updatedLines = append(updatedLines, "\terc20types \"github.com/cosmos/evm/x/erc20/types\"")
			updatedLines = append(updatedLines, "\tfeemarkettypes \"github.com/cosmos/evm/x/feemarket/types\"")
			updatedLines = append(updatedLines, "\tevmtypes \"github.com/cosmos/evm/x/vm/types\"")
			continue
		}

		// Add EVM module accounts
		if strings.Contains(line, "{Account: icatypes.ModuleName},") {
			updatedLines = append(updatedLines, line)
			updatedLines = append(updatedLines, "\t\t{Account: evmtypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},")
			updatedLines = append(updatedLines, "\t\t{Account: erc20types.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},")
			updatedLines = append(updatedLines, "\t\t{Account: feemarkettypes.ModuleName},")
			continue
		}

		// Add EVM modules to ordering
		if strings.Contains(line, "epochstypes.ModuleName,") && strings.Contains(lines[i+1], "// ibc modules") {
			updatedLines = append(updatedLines, line)
			updatedLines = append(updatedLines, "\t\t\t\t// cosmos evm modules")
			updatedLines = append(updatedLines, "\t\t\t\terc20types.ModuleName,")
			updatedLines = append(updatedLines, "\t\t\t\tfeemarkettypes.ModuleName,")
			updatedLines = append(updatedLines, "\t\t\t\tevmtypes.ModuleName,")
			continue
		}

		updatedLines = append(updatedLines, line)
	}

	return os.WriteFile(configPath, []byte(strings.Join(updatedLines, "\n")), 0644)
}
