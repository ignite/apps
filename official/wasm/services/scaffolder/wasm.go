package scaffolder

import (
	"context"

	"wasm/templates/wasm"

	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/pkg/errors"
)

const (
	wasmRepo = "github.com/CosmWasm/wasmd@v0.50.0"
)

type (
	// wasmOptions represents configuration for the message scaffolding.
	wasmOptions struct {
		home string
	}
	// WasmOption configures the message scaffolding.
	WasmOption func(*wasmOptions)
)

// newWasmOptions returns a wasmOptions with default options.
func newWasmOptions() wasmOptions {
	return wasmOptions{}
}

// WithHome provides a custom chain home path.
func WithHome(home string) WasmOption {
	return func(m *wasmOptions) {
		m.home = home
	}
}

// AddWasm add wasm support.
func (s Scaffolder) AddWasm(
	ctx context.Context,
	tracer *placeholder.Tracer,
	options ...WasmOption,
) (xgenny.SourceModification, error) {
	path := s.chain.AppPath()
	if hasWasm(path) {
		return xgenny.SourceModification{}, errors.Errorf("wasm integration already exist for path %s", path)
	}

	home, err := s.chain.Home()
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	// Create the options
	scaffoldingOpts := newWasmOptions()
	for _, apply := range options {
		apply(&scaffoldingOpts)
	}
	opts := &wasm.Options{
		AppName: s.chain.Name(),
		AppPath: path,
		Home:    home,
	}

	// Scaffold
	g, err := wasm.NewWasmGenerator(tracer, opts)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	sm, err := xgenny.RunWithValidation(tracer, g)
	if err != nil {
		return sm, err
	}
	return sm, finish(ctx, opts.AppPath)
}
