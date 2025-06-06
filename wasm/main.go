package main

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/wasm/cmd"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	m := &plugin.Manifest{
		Name:     "wasm",
		Commands: cmd.GetCommands(),
	}
	return m, nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	// Remove the three two elements "ignite" and "wasm" from OsArgs.
	args := c.OsArgs[2:]
	switch args[0] {
	case "add":
		return cmd.AddHandler(ctx, c, api)
	case "config":
		return cmd.ConfigHandler(ctx, c, api)
	default:
		return errors.Errorf("unknown command: %s", c.Path)
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
			"wasm": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
