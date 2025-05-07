package main

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/cmd"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "hermes",
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// Remove the three elements "ignite", "relayer" and "hermes" from OsArgs.
	args := c.OsArgs[3:]
	switch args[0] {
	case "configure":
		return cmd.ConfigureHandler(ctx, c)
	case "exec":
		return cmd.ExecuteHandler(ctx, c)
	case "start":
		return cmd.StartHandler(ctx, c)
	case "keys":
		switch args[1] {
		case "add":
			return cmd.KeysAddMnemonicHandler(ctx, c)
		case "file":
			return cmd.KeysAddFileHandler(ctx, c)
		case "list":
			return cmd.KeysListHandler(ctx, c)
		case "delete":
			return cmd.KeysDeleteHandler(ctx, c)
		default:
			return errors.Errorf("unknown keys command: %s", args[1])
		}
	default:
		return errors.Errorf("unknown command: %s", args[0])
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
			"hermes": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
