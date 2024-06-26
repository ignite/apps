package main

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/cmd"
)

var _ plugin.Interface = app{}

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	m := &plugin.Manifest{
		Name:     "hermes",
		Commands: cmd.GetCommands(),
	}
	return m, nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// Instead of a switch on c.Use, we run the root command like if
	// we were in a command line context. This implies to set os.Args
	// correctly.
	// Remove the first arg "ignite" from OsArgs because our hermes
	// command root is "relayer" not "ignite".
	args := c.OsArgs[1:]
	switch args[0] {
	case "configure":
		return cmd.ConfigureHandler(ctx, c)
	case "exec":
		return cmd.ExecuteHandler(ctx, c)
	case "start":
		return cmd.StartHandler(ctx, c)
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
			"hermes": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
