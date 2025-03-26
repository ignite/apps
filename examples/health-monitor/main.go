package main

import (
	"context"
	"strings"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/examples/health-monitor/cmd"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
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
		return errors.Errorf("failed to get chain info: %s", err)
	}

	switch args[0] {
	case "monitor":
		return cmd.ExecuteMonitor(ctx, c, chainInfo)
	default:
		return errors.Errorf("unknown command: %s", strings.Join(c.OsArgs, " "))
	}
}

func (app) ExecuteHookPre(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPost(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookCleanUp(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
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
