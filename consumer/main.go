package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(ctx context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "consumer",
	}, nil
}

func (a app) Execute(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("missing argument")
	}
	chain, err := api.GetChainInfo(ctx)
	if err != nil {
		return err
	}
	switch cmd.Args[0] {
	case "writeGenesis":
		return writeConsumerGenesis(chain)
	case "isInitialized":
		isInit, err := isInitialized(chain)
		fmt.Printf("%t", isInit)
		return err
	}
	return fmt.Errorf("invalid argument %q", cmd.Args[0])
}

func (app) ExecuteHookPre(ctx context.Context, h *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPost(ctx context.Context, h *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookCleanUp(ctx context.Context, h *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"consumer": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
