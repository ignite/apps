package main

import (
	"context"
	"os"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/apps/explorer/cmd"

	"github.com/ignite/cli/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	m := &plugin.Manifest{Name: "explorer"}
	m.ImportCobraCommand(cmd.NewExplorer(), "ignite")
	return m, nil
}

func (app) Execute(_ context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// Instead of a switch on c.Use, we run the root command like if
	// we were in a command line context. This implies to set os.Args
	// correctly.
	// Remove the first arg "ignite" from OSArgs because our explorer
	// command root is "explorer" not "ignite".
	os.Args = c.OsArgs[1:]
	return cmd.NewExplorer().Execute()
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
			"explorer": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
