package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"health-monitor/cmd"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(_ context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "health-monitor",
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	// Remove the first two elements "ignite" and "health-monitor" from OsArgs.
	args := c.OsArgs[2:]

	chainInfo, err := api.GetChainInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain info: %w", err)
	}

	switch args[0] {
	case "monitor":
		return cmd.ExecuteMonitor(ctx, c, chainInfo)
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
			"health-monitor": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
