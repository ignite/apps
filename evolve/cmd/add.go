package cmd

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/evolve/template"
)

const (
	statusScaffolding = "Scaffolding..."

	flagPath    = "path"
	flagStart   = "start"
	flagMigrate = "migrate"
)

func AddHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	flags := plugin.Flags(cmd.Flags)

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	appPath, err := flags.GetString(flagPath)
	if err != nil {
		return err
	}

	withStartCmd, err := flags.GetBool(flagStart)
	if err != nil {
		return err
	}

	migrateCometBFT, err := flags.GetBool(flagMigrate)
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

	g, err := template.NewEvolveGenerator(c, migrateCometBFT, withStartCmd)
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

	err = session.Printf("🎉 Evolve (ev-abci) added (`%[1]v`).\n", c.AppPath(), c.Name())

	if migrateCometBFT {
		err = errors.Join(session.Printf("\n"))
		err = errors.Join(err, session.Println("Additionally, evolve migration commands and modules successfully scaffolded!"))
		err = errors.Join(err, session.Printf("If %s is already live, check out the newly added evolve manager to prepare the chain for migration\n", c.Name()))
		err = errors.Join(err, session.Printf("Run `%s evolve-migrate` to migrate CometBFT state to the evolve state.\n", binaryName))
	}

	return err
}

// finish finalize the scaffolded code (formating, dependencies).
func finish(ctx context.Context, session *cliui.Session, path string) error {
	session.StopSpinner()
	session.StartSpinner("go mod tidy...")
	if err := gocmd.ModTidy(ctx, path); err != nil {
		return err
	}

	session.StopSpinner()
	session.StartSpinner("Formatting code...")
	if err := gocmd.Fmt(ctx, path); err != nil {
		return err
	}

	_ = gocmd.GoImports(ctx, path) // goimports installation could fail, so ignore the error

	return nil
}
