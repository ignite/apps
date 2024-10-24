package config

import (
	"os"
	"strings"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

type (
	// options represents configuration for the wasm config.
	options struct {
		simulationGasLimit uint64
		smartQueryGasLimit uint64
		memoryCacheSize    uint64
	}
	// Option configures the message scaffolding.
	Option func(*options)
)

// newOptions returns a wasmOptions with default options.
func newOptions() options {
	return options{
		simulationGasLimit: 0,
		smartQueryGasLimit: 3_000_000,
		memoryCacheSize:    100,
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

// AddWasm add wasm parameters to the chain TOML config.
func AddWasm(configPath string, options ...Option) error {
	// Create the options.
	opts := newOptions()
	for _, apply := range options {
		apply(&opts)
	}

	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	// Check if the wasm section already exist.
	if hasWasm(configPath) {
		return errors.Errorf("config file already have wasm %s", configPath)
	}

	// Add default configs.
	config := wasmtypes.DefaultWasmConfig()
	if opts.simulationGasLimit != 0 {
		config.SimulationGasLimit = &opts.simulationGasLimit
	}
	config.SmartQueryGasLimit = opts.smartQueryGasLimit
	config.MemoryCacheSize = uint32(opts.memoryCacheSize)

	// Save new configs to the TOML file.
	if _, err = f.WriteString(wasmtypes.ConfigTemplate(config)); err != nil {
		return err
	}

	return nil
}

// hasWasm check if the TOML config already have the wasm section.
func hasWasm(configPath string) bool {
	f, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	if strings.Contains(string(f), "[wasm]") {
		return true
	}
	return false
}
