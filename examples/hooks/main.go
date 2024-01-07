package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/apps/examples/hooks/cmd"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(_ context.Context) (*plugin.Manifest, error) {
	m := &plugin.Manifest{
		Name: "hooks",
		Hooks: []*plugin.Hook{
			{
				Name:        "chain-scaffold",
				PlaceHookOn: "ignite scaffold chain",
			},
			{
				Name:        "chain-serve",
				PlaceHookOn: "ignite chain serve",
			},
		},
	}
	m.ImportCobraCommand(cmd.NewHooks(), "ignite")

	return m, nil
}

func (app) Execute(_ context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPre(_ context.Context, hook *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	fmt.Printf("ExecuteHookPre: %s\n", hook.Hook.Name)
	return nil
}

func (app) ExecuteHookPost(_ context.Context, hook *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	fmt.Printf("ExecuteHookPost: %s\n", hook.Hook.Name)
	return nil
}

func (app) ExecuteHookCleanUp(_ context.Context, hook *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	fmt.Printf("ExecuteHookCleanUp: %s\n", hook.Hook.Name)
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"hooks": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
