package scaffolder

import (
	"context"
	"fmt"
	"os"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"

	"github.com/ignite/apps/wasm/pkg/config"
	"github.com/ignite/apps/wasm/pkg/xgit"
	"github.com/ignite/apps/wasm/templates/wasm"
)

const (
	wasmRepo = "github.com/CosmWasm/wasmd"

	errNotCompatibleVersionStr = `The WASM version (%s) is not compatible with the Cosmos SDK version (%s) of your chain`
	errOldCosmosSDKVersionStr  = `Your chain has been scaffolded with an older version of Cosmos SDK: %s

Please, follow the migration guide to upgrade your chain to the latest version at https://docs.ignite.com/migration`
)

var (
	wasmV601 = semver.MustParse("0.60.1")
	wasmV600 = semver.MustParse("0.60.0")
	wasmV551 = semver.MustParse("0.55.1")

	LegacyWasmVersion  = wasmV551
	DefaultWasmVersion = wasmV601
)

type (
	// options represent configuration for the message scaffolding.
	options struct {
		wasmVersion        semver.Version
		simulationGasLimit uint64
		smartQueryGasLimit uint64
		memoryCacheSize    uint64
	}
	// Option configures the message scaffolding.
	Option func(*options)
)

// newOptions returns a options with default options.
func newOptions() options {
	return options{
		simulationGasLimit: 0,
		smartQueryGasLimit: 3_000_000,
		memoryCacheSize:    100,
	}
}

// WithWasmVersion set the wasm semantic version.
func WithWasmVersion(version semver.Version) Option {
	return func(m *options) {
		m.wasmVersion = version
	}
}

// WithSimulationGasLimit provides a simulation gas limit for the wasm config.
func WithSimulationGasLimit(simulationGasLimit uint64) Option {
	return func(m *options) {
		m.simulationGasLimit = simulationGasLimit
	}
}

// WithSmartQueryGasLimit provides a smart query gas limit for the wasm config.
func WithSmartQueryGasLimit(smartQueryGasLimit uint64) Option {
	return func(m *options) {
		m.smartQueryGasLimit = smartQueryGasLimit
	}
}

// WithMemoryCacheSize provides a memory cache size for the wasm config.
func WithMemoryCacheSize(memoryCacheSize uint64) Option {
	return func(m *options) {
		m.memoryCacheSize = memoryCacheSize
	}
}

// AddWasm add wasm support.
func (s Scaffolder) AddWasm(
	ctx context.Context,
	tracer *placeholder.Tracer,
	options ...Option,
) (xgenny.SourceModification, error) {
	scaffoldingOpts := newOptions()
	for _, apply := range options {
		apply(&scaffoldingOpts)
	}

	// Check if the wasm version exist
	versions, err := xgit.FetchGitTags(fmt.Sprintf("https://%s", wasmRepo))
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	if !xgit.HasVersion(versions, scaffoldingOpts.wasmVersion) {
		return xgenny.SourceModification{},
			errors.Errorf("semantic version v%s not exist in %s", scaffoldingOpts.wasmVersion.String(), wasmRepo)
	}

	// Check if the chain already has wasm integration.
	path := s.chain.AppPath()
	if hasWasm(path) {
		return xgenny.SourceModification{}, errors.Errorf("wasm integration already exist for path %s", path)
	}

	// Prepare scaffold.
	home, err := s.chain.Home()
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	appTOMLPath, err := s.chain.AppTOMLPath()
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	if _, err := os.Stat(appTOMLPath); os.IsNotExist(err) {
		s.session.Printf(`Cannot find the chain config. If the chain %[1]v is not initialized yet, run "%[1]vd init" or "ignite chain serve" to init the chain. 
After, run the "ignite wasm config" command to add the wasm config

`,
			s.chain.Name(),
		)
	} else if err == nil {
		// Add wasm options to the chain config.
		if err := config.AddWasm(
			appTOMLPath,
			config.WithSimulationGasLimit(scaffoldingOpts.simulationGasLimit),
			config.WithSmartQueryGasLimit(scaffoldingOpts.smartQueryGasLimit),
			config.WithMemoryCacheSize(scaffoldingOpts.memoryCacheSize),
		); err != nil {
			return xgenny.SourceModification{}, err
		}
	}

	// Scaffold wasm changes.
	binaryName, err := s.chain.Binary()
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	// Scaffold wasm changes.
	legacyWasm, err := assertVersions(scaffoldingOpts.wasmVersion, s.chain.Version)
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	opts := &wasm.Options{
		BinaryName: binaryName,
		AppPath:    path,
		Home:       home,
		Legacy:     legacyWasm,
	}
	g, err := wasm.NewWasmGenerator(tracer, opts)
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	sm, err := xgenny.NewRunner(ctx, opts.AppPath).RunAndApply(g)
	if err != nil {
		return sm, err
	}

	return sm, finish(ctx, s.session, opts.AppPath, scaffoldingOpts.wasmVersion)
}

// assertVersions asserts that the given wasm and sdk versions are compatible.
func assertVersions(wasmVersion semver.Version, sdkVersion cosmosver.Version) (legacyWasm bool, err error) {
	switch {
	case sdkVersion.LT(cosmosver.StargateFiftyVersion):
		// Invalid SDK version should be at least Stargate 0.50
		return false, errors.Errorf(errOldCosmosSDKVersionStr, sdkVersion)
	case wasmVersion.LT(wasmV600):
		// if the wasm version is legacy, check if it is compatible with the SDK version
		if sdkVersion.LT(cosmosver.StargateFiftyThreeVersion) {
			return false, errors.Errorf(errNotCompatibleVersionStr, wasmVersion, sdkVersion)
		}
		return true, nil
	default:
		if sdkVersion.LT(cosmosver.StargateFiftyThreeVersion) {
			return false, errors.Errorf(errNotCompatibleVersionStr, wasmVersion, sdkVersion)
		}
		return false, nil
	}
}
