package cmd

import (
	"context"

	"github.com/ignite/cli/v29/ignite/services/plugin"
)

// EditGenesisHandler allows to edit the genesis file to add the sequencer module.
func EditGenesisHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	return initEVABCI(ctx, cmd, false)
}
