package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

var _ plugin.Interface = app{}

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "consumer",
	}, nil
}

func (a app) Execute(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	if len(cmd.Args) == 0 {
		return errors.Errorf("missing argument")
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
	return errors.Errorf("invalid argument %q", cmd.Args[0])
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
			"consumer": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
