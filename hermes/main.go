package main

import (
	"context"
	"os"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/apps/hermes/cmd"

	"github.com/ignite/cli/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	m := &plugin.Manifest{Name: "hermes"}
	m.ImportCobraCommand(cmd.NewHermes(), "ignite relayer")
	return m, nil
}

func (app) Execute(_ context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// Run the "hermes" command as if it were a root command. To do
	// so remove the first two arguments which are "ignite relayer"
	// from OSArgs to treat "hermes" as the root command.
	os.Args = c.OsArgs[2:]
	return cmd.NewHermes().Execute()
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
