package main

import (
	"context"
	"fmt"
	"os"

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
	// Instead of a switch on c.Use, we run the root command like if
	// we were in a command line context. This implies to set os.Args
	// correctly.
	// Remove the first arg "ignite" from OSArgs because our hooks
	// command root is "hooks" not "ignite".
	os.Args = c.OsArgs[1:]
	return cmd.NewHooks().Execute()
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

		fmt.Println("chain-id:", chainInfo.ChainId)
		fmt.Println("rpc-address:", chainInfo.RpcAddress)
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
