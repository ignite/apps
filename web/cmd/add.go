package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ignite/apps/web/templates"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

const (
	statusScaffolding = "Scaffolding..."

	flagPath = "path"
)

// ExecuteAdd executes the add subcommand.
func ExecuteAdd(ctx context.Context, cmd *plugin.ExecutedCommand) error {
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

	if err := templates.Write(c.AppPath()); err != nil {
		return fmt.Errorf("failed to write chain-admin: %w", err)
	}

	return session.Printf("🎉 Ignite chain-admin added (`%[1]v`).\n", c.AppPath(), c.Name())
}
