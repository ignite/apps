package cmd

import (
	"context"
	"fmt"

	"github.com/ignite/cli/v28/ignite/pkg/chainregistry"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/connect/chains"
)

const (
	fetchingChainInfo = "Fetching chain info..."
)

func AddHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s", cmd.Use)
	}

	session := cliui.New(cliui.StartSpinnerWithText(fetchingChainInfo))
	defer session.End()

	chainRegistry := chains.NewChainRegistry()
	if err := chainRegistry.FetchChains(); err != nil {
		return err
	}

	chain, ok := chainRegistry.Chains[cmd.Args[0]]
	if !ok {
		return fmt.Errorf("chain %s not found", cmd.Args[0])
	}

	if err := initializeChain(ctx, chain); err != nil {
		return err
	}

	return nil
}

func initializeChain(ctx context.Context, chain chainregistry.Chain) error {
	fmt.Println("Initializing chain", chain.ChainName)
	return nil
}
