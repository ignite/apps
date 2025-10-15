package cmd

import (
	"context"
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

	g, err := template.NewEvolveGenerator(c, template.GeneratorOptions{
		WithStart: true,
	})
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

	return session.Printf("🎉 Evolve (ev-abci) added (`%[1]v`).\n", c.AppPath())
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
