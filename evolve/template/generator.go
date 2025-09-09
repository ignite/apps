package template

import (
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

// NewEvolveGenerator returns the generator to scaffold a evolve integration inside an app.
func NewEvolveGenerator(chain *chain.Chain, withCometMigration bool) (*genny.Generator, error) {
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

	g.RunFn(commandsStartModify(appPath, binaryName, chain.Version))
	g.RunFn(commandsGenesisModify(appPath, binaryName))
	g.RunFn(commandsRollbackModify(appPath, binaryName))
	if withCometMigration {
		g.RunFn(migrateFromCometModify(appPath))
	}

	return g, nil
}
