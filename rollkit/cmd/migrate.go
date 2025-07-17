package cmd

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

func MigrateFromCometHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
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

	g := &genny.Generator{} // TODO

	_, err = xgenny.RunWithValidation(placeholder.New(), g)
	if err != nil {
		return err
	}

	if finish(ctx, session, c.AppPath()) != nil {
		return err
	}

	binaryName, err := c.Binary()
	if err != nil {
		return err
	}

	err = errors.Join(err, session.Println("Rollkit migration commands and modules successfully scaffolded 🎉"))
	err = errors.Join(err, session.Printf("If %s is already live, check out the newly added rollkit manager to prepare the chain for migration\n", c.Name()))
	err = errors.Join(err, session.Printf("Run `%s rollkit-migrate` to migrate CometBFT state to the rollkit state.\n", binaryName))

	return err
}
