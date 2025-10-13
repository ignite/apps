package app

import (
	"fmt"
	"hash/fnv"
	"maps"
	"os"
	"path/filepath"

	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmempool "github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cast"

	evmconfig "github.com/cosmos/evm/config"
	evmmempool "github.com/cosmos/evm/mempool"
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
	"github.com/ethereum/go-ethereum/common"
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

// setEVMMempool sets the EVM priority nonce mempool
// it is required for the ethereum json rpc server to work
func (app *App) setEVMMempool() {
	if evmtypes.GetChainConfig() != nil {
		mempoolConfig := &evmmempool.EVMMempoolConfig{
			AnteHandler:   app.BaseApp.AnteHandler(),
			BlockGasLimit: 100_000_000,
		}

		evmMempool := evmmempool.NewExperimentalEVMMempool(app.CreateQueryContext, app.Logger(), app.EVMKeeper, app.FeeMarketKeeper, app.txConfig, app.clientCtx, mempoolConfig)
		app.EVMMempool = evmMempool

		app.SetMempool(evmMempool)
		checkTxHandler := evmmempool.NewCheckTxHandler(evmMempool)
		app.SetCheckTxHandler(checkTxHandler)

		abciProposalHandler := baseapp.NewDefaultProposalHandler(evmMempool, app)
		abciProposalHandler.SetSignerExtractionAdapter(evmmempool.NewEthSignerExtractionAdapter(sdkmempool.NewDefaultSignerExtractionAdapter()))
		app.SetPrepareProposal(abciProposalHandler.PrepareProposalHandler())
	}
}

// RegisterPendingTxListener a function that registers a listener for pending transactions.
func (app *App) RegisterPendingTxListener(listener func(common.Hash)) {
	app.pendingTxListeners = append(app.pendingTxListeners, listener)
}

// SetClientCtx a function that sets the client context on the app, required by EVM module implementation.
func (app *App) SetClientCtx(ctx client.Context) {
	app.clientCtx = ctx
}

// GetMempool returns the mempool of the app.
// It is required by the EVM application interface.
func (app *App) GetMempool() sdkmempool.ExtMempool {
	return app.EVMMempool
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
