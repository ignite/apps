package cmd

import (
	"context"

	"github.com/ignite/apps/connect/chains"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

func AppHandler(ctx context.Context, cmd *plugin.ExecutedCommand, name string, cfg chains.Config) error {
	return nil
}
