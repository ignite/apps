package main

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/fee-abstraction/cmd"
)

type app struct{}

func (app) Manifest(_ context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "fee-abstraction",
		Hooks:    cmd.GetHooks(),
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(_ context.Context, _ *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPre(_ context.Context, h *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	return cmd.ExecuteScaffoldPreHook(h)
}

func (app) ExecuteHookPost(ctx context.Context, h *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	if h.Hook.Name == cmd.ScaffoldChainHook {
		return cmd.ExecuteScaffoldChainPostHook(ctx, h)
	}
	return nil
}

func (app) ExecuteHookCleanUp(_ context.Context, _ *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"fee-abstraction": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
