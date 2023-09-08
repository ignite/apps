package main

import (
	"encoding/gob"
	"os"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/ignite/services/plugin"

	"explorer/cmd"
)

func init() {
	gob.Register(plugin.Manifest{})
	gob.Register(plugin.ExecutedCommand{})
	gob.Register(plugin.ExecutedHook{})
}

var _ plugin.Interface = (*p)(nil)

type p struct{}

func (p) Manifest() (plugin.Manifest, error) {
	m := plugin.Manifest{
		Name: "explorer",
	}
	m.ImportCobraCommand(cmd.NewExplorer(), "ignite")
	return m, nil
}

func (p) Execute(c plugin.ExecutedCommand) error {
	// Instead of a switch on c.Use, we run the root command like if
	// we were in a command line context. This implies to set os.Args
	// correctly.
	// Remove the first arg "ignite" from OSArgs because our explorer
	// command root is "explorer" not "ignite".
	os.Args = c.OSArgs[1:]
	return cmd.NewExplorer().Execute()
}

func (p) ExecuteHookPre(plugin.ExecutedHook) error {
	return nil
}

func (p) ExecuteHookPost(plugin.ExecutedHook) error {
	return nil
}

func (p) ExecuteHookCleanUp(plugin.ExecutedHook) error {
	return nil
}

func main() {
	pluginMap := map[string]hplugin.Plugin{
		"explorer": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
