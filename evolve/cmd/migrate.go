package cmd

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/evolve/template"
)

func MigrateHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	flags := plugin.Flags(cmd.Flags)

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	appPath, err := flags.GetString(flagPath)
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return err
	}

	c, err := chain.New(absPath, chain.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	binaryName, err := c.Binary()
	if err != nil {
		return err
	}

	g, err := template.NewEvolveMigrateGenerator(c)
	if err != nil {
		return err
	}

	_, err = xgenny.NewRunner(ctx, appPath).RunAndApply(g)
	if err != nil {
		return err
	}

	if finish(ctx, session, c.AppPath()) != nil {
		return err
	}

	err = session.Printf("🎉 Evolve migration support added (`%[1]v`).\n", c.AppPath())
	err = errors.Join(err, session.Println("Evolve migration command and module successfully scaffolded!"))
	err = errors.Join(err, session.Printf("If %s is already live, check out the newly added evolve manager to prepare the chain for migration\n", c.Name()))
	err = errors.Join(err, session.Printf("Run `%s evolve-migrate` to migrate CometBFT state to the evolve state.\n", binaryName))

	return err
}
