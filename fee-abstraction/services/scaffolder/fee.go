package scaffolder

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"

	"github.com/ignite/apps/fee-abstraction/pkg/xgit"
	"github.com/ignite/apps/fee-abstraction/template"
)

const (
	feeAbsRepo = "github.com/osmosis-labs/fee-abstraction"
	feeAbsPkg  = "github.com/osmosis-labs/fee-abstraction/v8"
)

type (
	// options represents configuration for the message scaffolding.
	options struct {
		version semver.Version
	}
	// Option configures the message scaffolding.
	Option func(*options)
)

// newOptions returns a options with default options.
func newOptions() options {
	return options{}
}

// WithVersion set the fee abstraction semantic version.
func WithVersion(version semver.Version) Option {
	return func(m *options) {
		m.version = version
	}
}

// AddFeeAbstraction add fee abstraction support.
func (s Scaffolder) AddFeeAbstraction(
	ctx context.Context,
	options ...Option,
) (xgenny.SourceModification, error) {
	scaffoldingOpts := newOptions()
	for _, apply := range options {
		apply(&scaffoldingOpts)
	}

	// Check if the fee abstraction version exists
	versions, err := xgit.FetchGitTags(fmt.Sprintf("https://%s", feeAbsRepo))
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	if !xgit.HasVersion(versions, scaffoldingOpts.version) {
		return xgenny.SourceModification{},
			errors.Errorf("semantic version v%s not exist in %s", scaffoldingOpts.version.String(), feeAbsRepo)
	}

	// Prepare scaffold.
	home, err := s.chain.Home()
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	// Scaffold fee abstraction changes.
	binaryName, err := s.chain.Binary()
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	path, err := filepath.Abs(s.chain.AppPath())
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	opts := &template.Options{
		BinaryName: binaryName,
		AppPath:    path,
		Home:       home,
	}

	runner := xgenny.NewRunner(ctx, path)
	g, err := template.NewFeeAbstractionGenerator(runner.Tracer(), opts)
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	sm, err := runner.RunAndApply(g)
	if err != nil {
		return sm, err
	}

	return sm, finish(ctx, s.session, opts.AppPath, scaffoldingOpts.version)
}
