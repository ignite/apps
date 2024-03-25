package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"hooks/cmd"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "hooks",
		Hooks: []*plugin.Hook{
			{
				Name:        "chain-build",
				PlaceHookOn: "ignite chain build",
			},
			{
				Name:        "chain-serve",
				PlaceHookOn: "ignite chain serve",
			},
		},
		Commands: cmd.GetCommands(),
	}, nil
}

func (app) Execute(context.Context, *plugin.ExecutedCommand, plugin.ClientAPI) error {
	fmt.Println(`To use either run "ignite chain build" or "ignite chain serve" and see the output.`)
	return nil
}

// ExecuteHookPre is called before the hook is executed.
// You can access the arguments and flags of the command that triggered the hook via hook.Command
// and some information like chain-id and rpc address of the chain in the current directory that
// ignite is running (if any) via the api.GetChainInfo.
func (app) ExecuteHookPre(ctx context.Context, hook *plugin.ExecutedHook, api plugin.ClientAPI) error {
	fmt.Printf("ExecuteHookPre: %s\n", hook.Hook.Name)
	if api != nil {
		chainInfo, err := api.GetChainInfo(ctx)
		if err != nil {
			return err
		}

		fmt.Printf(`Chain with chain-id "%s" running at rpc address "%s"\n`, chainInfo.ChainId, chainInfo.RpcAddress)
	}
	return nil
}

// ExecuteHookPost is called after the hook command is executed successfully.
func (app) ExecuteHookPost(_ context.Context, hook *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	fmt.Printf("ExecuteHookPost: %s\n", hook.Hook.Name)
	return nil
}

// ExecuteHookCleanUp is called after the hook command is executed, regardless of the result.
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
