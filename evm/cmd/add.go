package cmd

import (
	"context"

	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const (
	statusScaffolding = "Scaffolding..."

	flagPath = "path"
)

func AddHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	return nil
}
