package template

import (
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

// NewEvolveAddGenerator returns the generator to scaffold evolve integration with start command support.
func NewEvolveAddGenerator(chain *chain.Chain) (*genny.Generator, error) {
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

	g.RunFn(appConfigStakingModify(appPath))
	g.RunFn(commandsStartModify(appPath, binaryName, chain.Version))
	g.RunFn(commandsGenesisInitModify(appPath, binaryName))
	g.RunFn(commandsRollbackModify(appPath, binaryName))

	return g, nil
}

// NewEvolveMigrateGenerator returns the generator to scaffold evolve migration support.
func NewEvolveMigrateGenerator(chain *chain.Chain) (*genny.Generator, error) {
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

	g.RunFn(appConfigMigrateModify(appPath))
	g.RunFn(appConfigStakingModify(appPath))
	g.RunFn(commandsMigrateModify(appPath, binaryName))

	return g, nil
}
