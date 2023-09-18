package main

import (
	"encoding/gob"
	"os"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/ignite/services/plugin"

	"github.com/ignite/apps/hermes/cmd"
)

func init() {
	gob.Register(plugin.Manifest{})
	gob.Register(plugin.ExecutedCommand{})
	gob.Register(plugin.ExecutedHook{})
}

type p struct{}

func (p) Manifest() (plugin.Manifest, error) {
	m := plugin.Manifest{
		Name: "hermes",
	}
	cobraCmd, err := cmd.NewRelayer()
	if err != nil {
		return m, err
	}
	m.ImportCobraCommand(cobraCmd, "ignite")
	return m, nil
}

func (p) Execute(c plugin.ExecutedCommand) error {
	// Instead of a switch on c.Use, we run the root command like if
	// we were in a command line context. This implies to set os.Args
	// correctly.
	// Remove the first arg "ignite" from OSArgs because our network
	// command root is "network" not "ignite".
	os.Args = c.OSArgs[1:]
	cobraCmd, err := cmd.NewRelayer()
	if err != nil {
		return err
	}
	return cobraCmd.Execute()
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
		"hermes": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
