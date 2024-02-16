package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"chain-info/cmd"
)

type app struct{}

func (app) Manifest(_ context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "chain-info",
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	// Remove the first two elements "ignite" and "flags" from OsArgs.
	args := c.OsArgs[2:]

	chainInfo, err := api.GetChainInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain info: %w", err)
	}

	ch, err := chain.New(chainInfo.AppPath)
	if err != nil {
		return fmt.Errorf("failed to create a new chain object from app path: %w", err)
	}

	switch args[0] {
	case "info":
		return cmd.ExecuteInfo(ctx, c, ch)
	case "build":
		return cmd.ExecuteBuild(ctx, c, ch)
	default:
		return fmt.Errorf("unknown command: %s", c.Path)
	}
}

func (app) ExecuteHookPre(_ context.Context, _ *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPost(_ context.Context, _ *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookCleanUp(_ context.Context, _ *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"chain-info": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
