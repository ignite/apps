package main

import (
	"context"
	"slices"
	"strings"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/apps/connect/chains"
	"github.com/ignite/apps/connect/cmd"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

var _ plugin.Interface = app{}

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	var availableChains []string
	if cfg, err := chains.ReadConfig(); err == nil {
		for name := range cfg.Chains {
			availableChains = append(availableChains, name)
		}
	}

	return &plugin.Manifest{
		Name:     "connect",
		Commands: cmd.GetCommands(availableChains),
	}, nil
}

func (app) Execute(ctx context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// Instead of a switch on c.Use, we run the root command like if
	// we were in a command line context. This implies to set os.Args
	// correctly.
	// Remove the first arg "ignite" from OSArgs because our connect
	// command root is "connect" not "ignite".
	args := c.OsArgs[2:]

	var availableChains []string
	cfg, err := chains.ReadConfig()
	if err == nil {
		for name := range cfg.Chains {
			availableChains = append(availableChains, name)
		}
	}

	switch args[0] {
	case "discover":
		return cmd.DiscoverHandler(ctx, c)
	case "add", "init":
		return cmd.AddHandler(ctx, c)
	case "remove", "rm":
		return cmd.RemoveHandler(ctx, c)
	case "version":
		return cmd.VersionHandler(ctx, c)
	default:
		if slices.Contains(availableChains, args[0]) {
			return cmd.AppHandler(ctx, c, args[0], *cfg)
		}

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
			"connect": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
