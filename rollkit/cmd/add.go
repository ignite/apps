package cmd

import (
	"context"
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/rollkit/template"
)

const (
	statusScaffolding = "Scaffolding..."

	flagPath = "path"
)

func AddHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
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

	g, err := template.NewRollKitGenerator(c)
	if err != nil {
		return err
	}

	_, err = xgenny.RunWithValidation(placeholder.New(), g)
	if err != nil {
		return err
	}

	if finish(ctx, session, c.AppPath()) != nil {
		return err
	}

	return session.Printf("🎉 RollKit added (`%[1]v`).\n", c.AppPath(), c.Name())
}

// finish finalize the scaffolded code (formating, dependencies)
func finish(ctx context.Context, session *cliui.Session, path string) error {
	session.StartSpinner("go mod tidy...")
	if err := gocmd.ModTidy(ctx, path); err != nil {
		return err
	}

	session.StartSpinner("Formatting code...")
	if err := gocmd.Fmt(ctx, path); err != nil {
		return err
	}

	_ = gocmd.GoImports(ctx, path) // goimports installation could fail, so ignore the error

	return nil
}
