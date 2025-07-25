package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"

	"github.com/ignite/apps/cca/templates"
)

const (
	statusScaffolding = "Scaffolding..."

	flagPath = "path"
)

// ExecuteScaffold executes the scaffold cca subcommand.
func ExecuteScaffold(ctx context.Context, cmd *plugin.ExecutedCommand) error {
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

	sc, err := scaffolder.New(ctx, absPath, "proto")
	if err != nil {
		return err
	}

	cfg, err := c.Config()
	if err != nil {
		return err
	}

	// add chain registry files
	// those are used for the wallet connector
	if err = sc.CreateChainRegistryFiles(c, cfg); err != nil {
		return err
	}

	// add cca files
	if err := templates.Write(c.AppPath()); err != nil {
		return fmt.Errorf("failed to write CCA: %w", err)
	}

	return session.Printf("🎉 Ignite CCA added (`%[1]v/web`).\n", c.AppPath(), c.Name())
}
