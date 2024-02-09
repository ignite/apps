package scaffolder

import (
	"context"

	"wasm/templates/wasm"

	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
)

const (
	wasmRepo = "github.com/CosmWasm/wasmd@v0.50.0"
)

type (
	// wasmOptions represents configuration for the message scaffolding.
	wasmOptions struct {
		description string
	}
	// WasmOption configures the message scaffolding.
	WasmOption func(*wasmOptions)
)

// newWasmOptions returns a wasmOptions with default options.
func newWasmOptions() wasmOptions {
	return wasmOptions{
		description: "test",
	}
}

// WithDescription provides a custom description for the message CLI command.
func WithDescription(desc string) WasmOption {
	return func(m *wasmOptions) {
		m.description = desc
	}
}

// AddWasm add wasm support.
func (s Scaffolder) AddWasm(
	ctx context.Context,
	tracer *placeholder.Tracer,
	options ...WasmOption,
) (sm xgenny.SourceModification, err error) {
	// Create the options
	scaffoldingOpts := newWasmOptions()
	for _, apply := range options {
		apply(&scaffoldingOpts)
	}

	opts := &wasm.Options{
		AppName: s.modpath.Package,
		AppPath: s.path,
	}

	// Scaffold
	g, err := wasm.NewWasmGenerator(tracer, opts)
	if err != nil {
		return sm, err
	}
	sm, err = xgenny.RunWithValidation(tracer, g)
	if err != nil {
		return sm, err
	}
	return sm, finish(ctx, opts.AppPath)
}
