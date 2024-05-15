package cmd

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/services/chain"

	"github.com/ignite/apps/rollkit/template"
)

const (
	statusScaffolding = "Scaffolding..."

	flagPath = "path"
)

func NewRollkitAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "Add rollkit support",
		Long:  "Add rollkit support to your Cosmos SDK chain",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
			defer session.End()

			appPath, err := cmd.Flags().GetString(flagPath)
			if err != nil {
				return err
			}
			absPath, err := filepath.Abs(appPath)
			if err != nil {
				return err
			}

			chain, err := chain.New(absPath)
			if err != nil {
				return err
			}

			g, err := template.NewRollKitGenerator(chain)
			if err != nil {
				return err
			}

			_, err = xgenny.RunWithValidation(placeholder.New(), g)
			if err != nil {
				return err
			}

			if finish(cmd.Context(), session, chain.AppPath()) != nil {
				return err
			}

			session.Printf("\n🎉 RollKit added (`%[1]v`).\n\n", chain.AppPath())
			return nil
		},
	}

	c.Flags().StringP(flagPath, "p", ".", "path of the app")

	return c
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
