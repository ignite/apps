package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/v28/ignite/services/plugin"
	"wasm/cmd"
)

type app struct{}

func (app) Manifest(_ context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "wasm",
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// Remove the first two elements "ignite" and "wasm" from OsArgs.
	args := c.OsArgs[2:]

	switch args[0] {
	case "hello":
		return cmd.ExecuteHello(ctx, c)
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
			"wasm": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
