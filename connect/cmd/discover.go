package cmd

import (
	"context"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/connect/chains"
)

const (
	discoveringChains = "Discovering chains..."
)

func DiscoverHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	session := cliui.New(cliui.StartSpinnerWithText(discoveringChains))
	defer session.End()

	chainRegistry := chains.NewChainRegistry()
	if err := chainRegistry.FetchChains(); err != nil {
		return err
	}

	return nil
}
