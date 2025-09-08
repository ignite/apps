package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const (
	statusScaffolding = "Scaffolding EVM integration..."
	statusGenerating  = "Generating code..."
	statusUpdating    = "Updating dependencies..."

	flagPath = "path"
)

// AddHandler implements the EVM integration command
func AddHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	var (
		appPath = flagGetPath(cmd)
	)

	// Parse the app path to get module information
	appPath, err := filepath.Abs(appPath)
	if err != nil {
		return err
	}

	// Check if this is a valid Cosmos SDK app
	if err := validateCosmosApp(appPath); err != nil {
		return err
	}

	// Get module path from go.mod
	modulePath, appName, err := getModuleInfo(appPath)
	if err != nil {
		return err
	}

	session.EventBus().SendInfo("Adding EVM support to " + appName)

	// Update dependencies first
	session.StartSpinner(statusUpdating)
	if err := updateGoMod(appPath); err != nil {
		return errors.Wrap(err, "failed to update go.mod")
	}

	// Generate EVM integration files
	session.StartSpinner(statusGenerating)
	if err := scaffoldEVMIntegration(appPath, modulePath, appName); err != nil {
		return errors.Wrap(err, "failed to scaffold EVM integration")
	}

	// Update existing files
	if err := updateExistingFiles(appPath, modulePath, appName); err != nil {
		return errors.Wrap(err, "failed to update existing files")
	}

	session.StopSpinner()
	session.Printf("🎉 EVM integration added successfully!\n\n")
	session.Printf("Next steps:\n")
	session.Printf("1. Run `go mod tidy` to clean up dependencies\n")
	session.Printf("2. Build your application: `go build ./cmd/%sd`\n", appName)
	session.Printf("3. Initialize your chain with EVM support enabled\n")
	session.Printf("4. Start your node with `--evm.tracer=` flag for EVM tracing (optional)\n\n")

	return nil
}

func flagGetPath(cmd *plugin.ExecutedCommand) string {
	for _, flag := range cmd.Flags {
		if flag.Name == flagPath {
			if flag.Value != "" {
				return flag.Value
			}
			break
		}
	}
	return "."
}

func validateCosmosApp(appPath string) error {
	// Check for go.mod
	goModPath := filepath.Join(appPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return errors.New("no go.mod found - not a valid Go module")
	}

	// Check for app directory structure
	appDir := filepath.Join(appPath, "app")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		return errors.New("no app/ directory found - not a valid Cosmos SDK app")
	}

	return nil
}

func getModuleInfo(appPath string) (modulePath, appName string, err error) {
	goModPath := filepath.Join(appPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			modulePath = strings.TrimSpace(strings.TrimPrefix(line, "module"))
			parts := strings.Split(modulePath, "/")
			appName = parts[len(parts)-1]
			return modulePath, appName, nil
		}
	}

	return "", "", errors.New("module declaration not found in go.mod")
}

func updateGoMod(appPath string) error {
	goModPath := filepath.Join(appPath, "go.mod")

	// Read current go.mod
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string

	// Find module declaration and add replace directive after it
	replaceAdded := false
	requireStarted := false

	for i, line := range lines {
		updatedLines = append(updatedLines, line)

		// Add replace directive after module declaration
		if !replaceAdded && strings.HasPrefix(line, "module ") {
			// Look ahead for replace or require block
			nextNonEmptyIdx := findNextNonEmptyLine(lines, i+1)
			if nextNonEmptyIdx == -1 || (!strings.HasPrefix(lines[nextNonEmptyIdx], "replace") && !strings.HasPrefix(lines[nextNonEmptyIdx], "require")) {
				updatedLines = append(updatedLines, "")
				updatedLines = append(updatedLines, "replace github.com/ethereum/go-ethereum => github.com/cosmos/go-ethereum v1.16.2-cosmos-1")
				updatedLines = append(updatedLines, "")
				replaceAdded = true
			}
		}

		// Add replace directive before require block if not already added
		if !replaceAdded && strings.HasPrefix(line, "replace (") {
			replaceAdded = true // Already has replace block
		}

		if !replaceAdded && strings.HasPrefix(line, "require (") {
			// Insert replace before require
			updatedLines = updatedLines[:len(updatedLines)-1] // Remove current require line
			updatedLines = append(updatedLines, "replace github.com/ethereum/go-ethereum => github.com/cosmos/go-ethereum v1.16.2-cosmos-1")
			updatedLines = append(updatedLines, "")
			updatedLines = append(updatedLines, line) // Add back require line
			replaceAdded = true
			requireStarted = true
			continue
		}

		if strings.HasPrefix(line, "require (") {
			requireStarted = true
		}

		// Add EVM dependencies in require block
		if requireStarted && strings.Contains(line, "github.com/cosmos/cosmos-sdk") {
			// Add EVM dependencies after cosmos-sdk
			updatedLines = append(updatedLines, "\tgithub.com/cosmos/evm v1.0.0-rc2.0.20250822211227-2d3df2ba510c")
			updatedLines = append(updatedLines, "\tgithub.com/ethereum/go-ethereum v1.15.11")
			updatedLines = append(updatedLines, "\tgithub.com/spf13/cast v1.9.2")
		}
	}

	// Write updated go.mod
	return os.WriteFile(goModPath, []byte(strings.Join(updatedLines, "\n")), 0644)
}

func findNextNonEmptyLine(lines []string, start int) int {
	for i := start; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "" {
			return i
		}
	}
	return -1
}

func scaffoldEVMIntegration(appPath, modulePath, appName string) error {
	// Create EVM-related files
	files := map[string]string{
		"app/ante.go":                generateAnteFile(modulePath, appName),
		"app/ante/ante.go":           generateAnteHandlerFile(),
		"app/ante/cosmos_handler.go": generateCosmosHandlerFile(),
		"app/ante/evm_handler.go":    generateEVMHandlerFile(),
		"app/evm.go":                 generateEVMFile(appName),
	}

	for filePath, content := range files {
		fullPath := filepath.Join(appPath, filePath)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func updateExistingFiles(appPath, modulePath, appName string) error {
	// Update app.go
	if err := updateAppGo(appPath, modulePath, appName); err != nil {
		return err
	}

	// Update app_config.go
	if err := updateAppConfigGo(appPath); err != nil {
		return err
	}

	// Update cmd/root.go
	if err := updateRootCmd(appPath, appName); err != nil {
		return err
	}

	// Update cmd/commands.go
	if err := updateCommands(appPath, appName); err != nil {
		return err
	}

	return nil
}

func updateAppGo(appPath, modulePath, appName string) error {
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
			updatedLines = append(updatedLines, "\t_ \"github.com/ethereum/go-ethereum/eth/tracers/js\"")
			updatedLines = append(updatedLines, "\t_ \"github.com/ethereum/go-ethereum/eth/tracers/native\"")
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

func updateRootCmd(appPath, appName string) error {
	rootPath := filepath.Join(appPath, "cmd", appName+"d", "cmd", "root.go")

	content, err := os.ReadFile(rootPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string

	for _, line := range lines {
		// Add EVM keyring import
		if strings.Contains(line, "\"github.com/cosmos/cosmos-sdk/x/auth/types\"") {
			updatedLines = append(updatedLines, line)
			updatedLines = append(updatedLines, "\tcosmosevmkeyring \"github.com/cosmos/evm/crypto/keyring\"")
			continue
		}

		// Update client context configuration
		if strings.Contains(line, "WithViper(app.Name)") {
			updatedLines = append(updatedLines, line+" // env variable prefix")
			updatedLines = append(updatedLines, "\t\tWithKeyringOptions(cosmosevmkeyring.Option()). // evm keyring capabilities")
			updatedLines = append(updatedLines, "\t\tWithLedgerHasProtobuf(true)")
			continue
		}

		updatedLines = append(updatedLines, line)
	}

	return os.WriteFile(rootPath, []byte(strings.Join(updatedLines, "\n")), 0644)
}

func updateCommands(appPath, appName string) error {
	commandsPath := filepath.Join(appPath, "cmd", appName+"d", "cmd", "commands.go")

	content, err := os.ReadFile(commandsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string

	for _, line := range lines {
		// Replace keys command with EVM-compatible version
		line = strings.Replace(line, "\"github.com/cosmos/cosmos-sdk/client/keys\"", "cosmosevmcmd \"github.com/cosmos/evm/client\"", 1)
		line = strings.Replace(line, "keys.Commands(),", "cosmosevmcmd.KeyCommands(app.DefaultNodeHome, true),", 1)

		updatedLines = append(updatedLines, line)
	}

	return os.WriteFile(commandsPath, []byte(strings.Join(updatedLines, "\n")), 0644)
}

// File generation functions
func generateAnteFile(modulePath, appName string) string {
	return fmt.Sprintf(`package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/evm/ante"
	evmante "github.com/cosmos/evm/ante"
	cosmosevmante "github.com/cosmos/evm/ante/evm"
	cosmosevmtypes "github.com/cosmos/evm/types"
	"github.com/ethereum/go-ethereum/common"

	appante "%s/app/ante"
)

// setAnteHandler sets the ante handler for the application.
func (app *App) setAnteHandler(txConfig client.TxConfig, maxGasWanted uint64) {
	options := ante.HandlerOptions{
		Cdc:                    app.appCodec,
		AccountKeeper:          app.AuthKeeper,
		BankKeeper:             app.BankKeeper,
		ExtensionOptionChecker: cosmosevmtypes.HasDynamicFeeExtensionOption,
		EvmKeeper:              app.EVMKeeper,
		FeegrantKeeper:         app.FeeGrantKeeper,
		IBCKeeper:              app.IBCKeeper,
		FeeMarketKeeper:        app.FeeMarketKeeper,
		SignModeHandler:        txConfig.SignModeHandler(),
		SigGasConsumer:         evmante.SigVerificationGasConsumer,
		MaxTxGasWanted:         maxGasWanted,
		TxFeeChecker:           cosmosevmante.NewDynamicFeeChecker(app.FeeMarketKeeper),
		PendingTxListener:      func(hash common.Hash) {},
	}
	if err := options.Validate(); err != nil {
		panic(err)
	}

	app.SetAnteHandler(appante.NewAnteHandler(options))
}
`, modulePath)
}

func generateAnteHandlerFile() string {
	return `package ante

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/evm/ante"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// NewAnteHandler returns an ante handler responsible for attempting to route an
// Ethereum or SDK transaction to an internal ante handler for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
func NewAnteHandler(options ante.HandlerOptions) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx)
		if ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 0 {
				switch typeURL := opts[0].GetTypeUrl(); typeURL {
				case "/cosmos.evm.vm.v1.ExtensionOptionsEthereumTx":
					// handle as *evmtypes.MsgEthereumTx
					anteHandler = newMonoEVMAnteHandler(options)
				case "/cosmos.evm.types.v1.ExtensionOptionDynamicFeeTx":
					// cosmos-sdk tx with dynamic fee extension
					anteHandler = newCosmosAnteHandler(options)
				default:
					return ctx, errorsmod.Wrapf(
						errortypes.ErrUnknownExtensionOptions,
						"rejecting tx with unsupported extension option: %s", typeURL,
					)
				}

				return anteHandler(ctx, tx, sim)
			}
		}

		// handle as totally normal Cosmos SDK tx
		switch tx.(type) {
		case sdk.Tx:
			anteHandler = newCosmosAnteHandler(options)
		default:
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}
}
`
}

func generateCosmosHandlerFile() string {
	return `package ante

import (
	baseevmante "github.com/cosmos/evm/ante"
	cosmosante "github.com/cosmos/evm/ante/cosmos"
	evmante "github.com/cosmos/evm/ante/evm"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ibcante "github.com/cosmos/ibc-go/v10/modules/core/ante"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	sdkvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// newCosmosAnteHandler creates the default ante handler for Cosmos transactions
func newCosmosAnteHandler(options baseevmante.HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		cosmosante.NewRejectMessagesDecorator(), // reject MsgEthereumTxs
		cosmosante.NewAuthzLimiterDecorator( // disable the Msg types that cannot be included on an authz.MsgExec msgs field
			sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
			sdk.MsgTypeURL(&sdkvesting.MsgCreateVestingAccount{}),
		),
		ante.NewSetUpContextDecorator(),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		cosmosante.NewMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		evmante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
	)
}
`
}

func generateEVMHandlerFile() string {
	return `package ante

import (
	"github.com/cosmos/evm/ante"
	evmante "github.com/cosmos/evm/ante/evm"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// newMonoEVMAnteHandler creates the sdk.AnteHandler implementation for the EVM transactions.
func newMonoEVMAnteHandler(options ante.HandlerOptions) sdk.AnteHandler {
	decorators := []sdk.AnteDecorator{
		evmante.NewEVMMonoDecorator(
			options.AccountKeeper,
			options.FeeMarketKeeper,
			options.EvmKeeper,
			options.MaxTxGasWanted,
		),
		ante.NewTxListenerDecorator(options.PendingTxListener),
	}

	return sdk.ChainAnteDecorators(decorators...)
}
`
}

func generateEVMFile(appName string) string {
	return `package app

import (
	"fmt"
	"hash/fnv"
	"maps"
	"os"
	"path/filepath"

	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cast"

	evmconfig "github.com/cosmos/evm/config"
	"github.com/cosmos/evm/precompiles/bech32"
	"github.com/cosmos/evm/precompiles/p256"
	srvflags "github.com/cosmos/evm/server/flags"
	erc20 "github.com/cosmos/evm/x/erc20"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/cosmos/evm/x/feemarket"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	"github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	gethvm "github.com/ethereum/go-ethereum/core/vm"
)

// registerEVMModules register EVM keepers and non dependency inject modules.
func (app *App) registerEVMModules(appOpts servertypes.AppOptions) error {
	// chain config
	chainID := getEVMChainID(appOpts)
	coinInfoMap := map[uint64]evmtypes.EvmCoinInfo{
		chainID: evmtypes.EvmCoinInfo{
			Denom:         sdk.DefaultBondDenom,
			ExtendedDenom: sdk.DefaultBondDenom,
			DisplayDenom:  sdk.DefaultBondDenom,
			Decimals:      evmtypes.SixDecimals, // in line with Cosmos SDK default decimals
		},
	}

	// configure evm modules
	if err := evmconfig.EvmAppOptionsWithConfig(
		chainID,
		coinInfoMap,
		getCustomEVMActivators(),
	); err != nil {
		return err
	}

	// set up non depinject support modules store keys
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(evmtypes.StoreKey),
		storetypes.NewKVStoreKey(feemarkettypes.StoreKey),
		storetypes.NewKVStoreKey(erc20types.StoreKey),
		storetypes.NewTransientStoreKey(evmtypes.TransientKey),
		storetypes.NewTransientStoreKey(feemarkettypes.TransientKey),
	); err != nil {
		return err
	}

	// set up EVM keeper
	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		app.appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.GetKey(feemarkettypes.StoreKey),
		app.UnsafeFindStoreKey(feemarkettypes.TransientKey),
	)

	// NOTE: it's required to set up the EVM keeper before the ERC-20 keeper, because it is used in its instantiation.
	app.EVMKeeper = evmkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(evmtypes.StoreKey),
		app.UnsafeFindStoreKey(evmtypes.TransientKey),
		app.GetStoreKeysMap(),
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AuthKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.FeeMarketKeeper,
		&app.ConsensusParamsKeeper,
		&app.Erc20Keeper,
		tracer,
	)

	app.Erc20Keeper = erc20keeper.NewKeeper(
		app.GetKey(erc20types.StoreKey),
		app.appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AuthKeeper,
		app.BankKeeper,
		app.EVMKeeper,
		app.StakingKeeper,
		&app.TransferKeeper,
	)

	// register evm modules
	if err := app.RegisterModules(
		vm.NewAppModule(app.EVMKeeper, app.AuthKeeper, app.AuthKeeper.AddressCodec()),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		erc20.NewAppModule(app.Erc20Keeper, app.AuthKeeper),
	); err != nil {
		return err
	}

	return nil
}

func (app *App) postRegisterEVMModules() error {
	// register precompiles on EVMKeeper
	const bech32PrecompileBaseGas = 6_000

	// secp256r1 precompile as per EIP-7212
	p256Precompile := &p256.Precompile{}

	bech32Precompile, err := bech32.NewPrecompile(bech32PrecompileBaseGas)
	if err != nil {
		return fmt.Errorf("failed to instantiate bech32 precompile: %w", err)
	}

	precompiles := maps.Clone(gethvm.PrecompiledContractsPrague) // clone from latest vm fork.
	precompiles[bech32Precompile.Address()] = bech32Precompile
	precompiles[p256Precompile.Address()] = p256Precompile

	// add more stateful precompiles here, if needed.

	_ = app.EVMKeeper.WithStaticPrecompiles(precompiles)
	return nil
}

// getCustomEVMActivators defines a map of opcode modifiers associated
// with a key defining the corresponding EIP.
func getCustomEVMActivators() map[int]func(*gethvm.JumpTable) {
	var (
		multiplier        = uint64(10)
		sstoreConstantGas = uint64(500)
	)

	return map[int]func(*gethvm.JumpTable){
		0o000: func(jt *gethvm.JumpTable) {
			// enable0000 contains the logic to modify the CREATE and CREATE2 opcodes
			// constant gas value.
			currentValCreate := jt[gethvm.CREATE].GetConstantGas()
			jt[gethvm.CREATE].SetConstantGas(currentValCreate * multiplier)

			currentValCreate2 := jt[gethvm.CREATE2].GetConstantGas()
			jt[gethvm.CREATE2].SetConstantGas(currentValCreate2 * multiplier)
		},
		0o001: func(jt *gethvm.JumpTable) {
			// enable0001 contains the logic to modify the CALL opcode
			// constant gas value.
			currentVal := jt[gethvm.CALL].GetConstantGas()
			jt[gethvm.CALL].SetConstantGas(currentVal * multiplier)
		},
		0o002: func(jt *gethvm.JumpTable) {
			// enable0002 contains the logic to modify the SSTORE opcode
			// constant gas value.
			jt[gethvm.SSTORE].SetConstantGas(sstoreConstantGas)
		},
	}
}

// getEVMChainID returns the EVM chain ID from the app options.
func getEVMChainID(appOpts servertypes.AppOptions) uint64 {
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		// fallback to genesis chain-id
		genesisPathCfg, _ := appOpts.Get("genesis_file").(string)
		if genesisPathCfg == "" {
			genesisPathCfg = filepath.Join("config", "genesis.json")
		}

		reader, err := os.Open(filepath.Join(DefaultNodeHome, genesisPathCfg))
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		chainID, err = genutiltypes.ParseChainIDFromGenesis(reader)
		if err != nil {
			panic(fmt.Errorf("failed to parse chain-id from genesis file: %w", err))
		}
	}

	return cosmosChainIDToEVMChainID(chainID)
}

// cosmosChainIDToEVMChainID converts a Cosmos chain ID to an EVM chain ID.
// This is an opinionated function to simplify chain id management.
// In theory, cosmos chain id and evm chain id are independent and can be managed separately.
func cosmosChainIDToEVMChainID(chainID string) uint64 {
	hasher := fnv.New32a()
	hasher.Write([]byte(chainID))
	return uint64(hasher.Sum32())
}

// RegisterEVM Since the EVM modules don't support dependency injection,
// we need to manually register the modules on the client side.
// This needs to be removed after EVM supports App Wiring.
func RegisterEVM(cdc codec.Codec, interfaceRegistry codectypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		evmtypes.ModuleName:       vm.NewAppModule(nil, authkeeper.AccountKeeper{}, interfaceRegistry.SigningContext().AddressCodec()),
		erc20types.ModuleName:     erc20.NewAppModule(erc20keeper.Keeper{}, authkeeper.AccountKeeper{}),
		feemarkettypes.ModuleName: feemarket.NewAppModule(feemarketkeeper.Keeper{}),
	}

	for _, m := range modules {
		if mr, ok := m.(module.AppModuleBasic); ok {
			mr.RegisterInterfaces(cdc.InterfaceRegistry())
		}
	}

	return modules
}

// ProvideMsgEthereumTxCustomGetSigner provides a custom signer for the MsgEthereumTx message.
func ProvideMsgEthereumTxCustomGetSigner() signing.CustomGetSigner {
	return evmtypes.MsgEthereumTxCustomGetSigner
}

// GetStoreKeysMap returns a map of store keys.
func (app *App) GetStoreKeysMap() map[string]*storetypes.KVStoreKey {
	storeKeysMap := make(map[string]*storetypes.KVStoreKey)
	for _, storeKey := range app.GetStoreKeys() {
		kvStoreKey, ok := app.UnsafeFindStoreKey(storeKey.Name()).(*storetypes.KVStoreKey)
		if ok {
			storeKeysMap[storeKey.Name()] = kvStoreKey
		}
	}

	return storeKeysMap
}
`
}
