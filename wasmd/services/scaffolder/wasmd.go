package scaffolder

import (
	"context"
	"errors"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/apps/wasmd/templates/initialize"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/validation"
	"github.com/ignite/cli/ignite/pkg/xgenny"
)

// wasmdOptions holds options for creating a new module.
type wasmdOptions struct {
	// ibc true if the module is an ibc module
	ibc bool

	// params list of parameters
	params []string
}

// WasmdOption configures Chain.
type WasmdOption func(*wasmdOptions)

// WithIBC scaffolds a module with IBC enabled.
func WithIBC() WasmdOption {
	return func(m *wasmdOptions) {
		m.ibc = true
	}
}

// WithParams scaffolds a module with params.
func WithParams(params []string) WasmdOption {
	return func(m *wasmdOptions) {
		m.params = params
	}
}

// CreateModule creates a new empty module in the scaffolded app.
func (s Scaffolder) InitWasmd(
	ctx context.Context,
	tracer *placeholder.Tracer,
	options ...WasmdOption,
) (sm xgenny.SourceModification, err error) {

	// Apply the options
	var wasmdOpts wasmdOptions
	for _, apply := range options {
		apply(&wasmdOpts)
	}

	opts := &initialize.InitOptions{
		AppName: s.modpath.Package,
		AppPath: s.path,
		Version: "v0.44",
	}

	g, err := initialize.NewGenerator(opts)
	if err != nil {
		return sm, err
	}

	gens := []*genny.Generator{g}

	sm, err = xgenny.RunWithValidation(tracer, gens...)
	if err != nil {
		return sm, err
	}

	// Modify app.go to register the module
	newSourceModification, runErr := xgenny.RunWithValidation(tracer, initialize.NewAppModify(tracer, opts))
	sm.Merge(newSourceModification)
	var validationErr validation.Error
	if runErr != nil && !errors.As(runErr, &validationErr) {
		return sm, runErr
	}
	if err := s.installWasm("v0.44.0"); err != nil {
		return sm, err
	}
	return sm, finish(ctx, opts.AppPath, s.modpath.RawPath)
}
