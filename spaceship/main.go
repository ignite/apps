package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/cmd"
)

type app struct{}

func (app) Manifest(_ context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "spaceship",
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	chainInfo, err := api.GetChainInfo(ctx)
	if err != nil {
		return errors.Errorf("failed to get chain info: %s", err)
	}

	// Remove the first two elements "ignite" and "spaceship" from OsArgs.
	args := c.OsArgs[2:]
	switch args[0] {
	case "deploy":
		return cmd.ExecuteSSHDeploy(ctx, c, chainInfo)
	case "log":
		return cmd.ExecuteSSHLog(ctx, c, chainInfo)
	case "status":
		return cmd.ExecuteSSHStatus(ctx, c, chainInfo)
	case "restart":
		return cmd.ExecuteSSHRestart(ctx, c, chainInfo)
	case "stop":
		return cmd.ExecuteSSHSStop(ctx, c, chainInfo)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
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
			"spaceship": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
