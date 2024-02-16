package scaffolder

import (
	"context"
	"os"

	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/pkg/errors"

	"wasm/pkg/config"
	"wasm/templates/wasm"
)

const (
	wasmRepo = "github.com/CosmWasm/wasmd@v0.50.0"
)

type (
	// wasmOptions represents configuration for the message scaffolding.
	wasmOptions struct {
		simulationGasLimit uint64
		smartQueryGasLimit uint64
		memoryCacheSize    uint32
	}
	// WasmOption configures the message scaffolding.
	WasmOption func(*wasmOptions)
)

// newWasmOptions returns a wasmOptions with default options.
func newWasmOptions() wasmOptions {
	return wasmOptions{
		simulationGasLimit: 0,
		smartQueryGasLimit: 3_000_000,
		memoryCacheSize:    100,
	}
}

// WithSimulationGasLimit provides a simulation gas limit for the wasm config.
func WithSimulationGasLimit(simulationGasLimit uint64) WasmOption {
	return func(m *wasmOptions) {
		m.simulationGasLimit = simulationGasLimit
	}
}

// WithSmartQueryGasLimit provides a smart query gas limit for the wasm config.
func WithSmartQueryGasLimit(smartQueryGasLimit uint64) WasmOption {
	return func(m *wasmOptions) {
		m.smartQueryGasLimit = smartQueryGasLimit
	}
}

// WithMemoryCacheSize provides a memory cache size for the wasm config.
func WithMemoryCacheSize(memoryCacheSize uint32) WasmOption {
	return func(m *wasmOptions) {
		m.memoryCacheSize = memoryCacheSize
	}
}

// AddWasm add wasm support.
func (s Scaffolder) AddWasm(
	ctx context.Context,
	tracer *placeholder.Tracer,
	options ...WasmOption,
) (xgenny.SourceModification, error) {
	// Check if chain already have wasm integration.
	path := s.chain.AppPath()
	if hasWasm(path) {
		return xgenny.SourceModification{}, errors.Errorf("wasm integration already exist for path %s", path)
	}

	// Prepare scaffold.
	home, err := s.chain.Home()
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	scaffoldingOpts := newWasmOptions()
	for _, apply := range options {
		apply(&scaffoldingOpts)
	}

	configTOML, err := s.chain.ConfigTOMLPath()
	if _, err := os.Stat(configTOML); os.IsNotExist(err) {
		s.session.Printf(`Cannot find the chain config. If the chain %[1]v is not initialized yet, run "%[1]vd init" or "ignite chain serve" to init the chain. 
After, run the "ignite wasm config" command to add the wasm config

`,
			s.chain.Name(),
		)
	} else if err == nil {
		// Add wasm options to the chain config.
		if err := config.AddWasm(
			configTOML,
			config.WithSimulationGasLimit(scaffoldingOpts.simulationGasLimit),
			config.WithSmartQueryGasLimit(scaffoldingOpts.smartQueryGasLimit),
			config.WithMemoryCacheSize(scaffoldingOpts.memoryCacheSize),
		); err != nil {
			return xgenny.SourceModification{}, err
		}
	}

	// Scaffold wasm changes.
	opts := &wasm.Options{
		AppName: s.chain.Name(),
		AppPath: path,
		Home:    home,
	}
	g, err := wasm.NewWasmGenerator(tracer, opts)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	sm, err := xgenny.RunWithValidation(tracer, g)
	if err != nil {
		return sm, err
	}

	return sm, finish(ctx, s.session, opts.AppPath)
}
