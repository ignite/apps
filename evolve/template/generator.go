package template

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

// GeneratorOptions represents the options for the generator.
type GeneratorOptions struct {
	WithStart     bool
	WithMigration bool
}

// NewEvolveGenerator returns the generator to scaffold evolve integration.
func NewEvolveGenerator(chain *chain.Chain, opts GeneratorOptions) (*genny.Generator, error) {
	g := genny.New()
	ctx := plush.NewContext()
	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))

	binaryName, err := chain.Binary()
	if err != nil {
		return nil, err
	}

	appPath := chain.AppPath()

	if err := updateDependencies(appPath); err != nil {
		return nil, errors.Errorf("failed to update go.mod: %w", err)
	}

	g.RunFn(appModify(appPath, opts.WithMigration))

	if opts.WithStart {
		g.RunFn(commandsStartModify(appPath, binaryName, chain.Version))
		g.RunFn(commandsGenesisInitModify(appPath, binaryName))
		g.RunFn(commandsRollbackModify(appPath, binaryName))
		g.RunFn(commandsForceInclusionModify(appPath, binaryName))
	}

	if opts.WithMigration {
		g.RunFn(commandsMigrateModify(appPath, binaryName))
	}

	return g, nil
}

// updateDependencies makes sure the correct dependencies are added to the go.mod files.
// ev-abci expects evolve v1 to be used.
func updateDependencies(appPath string) error {
	gomod, err := gomodule.ParseAt(appPath)
	if err != nil {
		return errors.Errorf("failed to parse go.mod: %w", err)
	}

	gomod.AddNewRequire(EvABCIPackage, EvABCIVersion, false)
	gomod.AddNewRequire(EvNodePackage, EvNodeVersion, false)

	// add local-da as go tool dependency (useful for local development)
	if err := gomod.AddTool(EvNodeDaCmd); err != nil {
		return errors.Errorf("failed to add local-da tool: %w", err)
	}

	// add required replaces
	gomod.AddReplace(GoHeaderPackage, "", GoHeaderPackageFork, GoHeaderVersionFork)

	// add temporary replaces
	// TODO(@julienrbrt): remove after tagged version of ev-abci and ev-node
	gomod.AddReplace("github.com/evstack/ev-node/core", "", "github.com/evstack/ev-node/core", "v1.0.0-beta.5.0.20251216132820-afcd6bd9b354")

	// save go.mod
	data, err := gomod.Format()
	if err != nil {
		return errors.Errorf("failed to format go.mod: %w", err)
	}

	return os.WriteFile(filepath.Join(appPath, "go.mod"), data, 0o644)
}
