package main

import (
	"encoding/gob"
	"os"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/ignite/services/plugin"

	"github.com/ignite/apps/airdrop/cmd"
)

func init() {
	gob.Register(plugin.Manifest{})
	gob.Register(plugin.ExecutedCommand{})
	gob.Register(plugin.ExecutedHook{})
}

type p struct{}

func (p) Manifest() (plugin.Manifest, error) {
	m := plugin.Manifest{
		Name: "airdrop",
	}
	m.ImportCobraCommand(cmd.NewAirdrop(), "ignite")
	return m, nil
}

func (p) Execute(c plugin.ExecutedCommand) error {
	// Remove the first arg "ignite" from OSArgs because our command root is
	// "airdrop" not "ignite".
	os.Args = c.OSArgs[1:]
	return cmd.NewAirdrop().Execute()
}

func (p) ExecuteHookPre(hook plugin.ExecutedHook) error {
	return nil
}

func (p) ExecuteHookPost(hook plugin.ExecutedHook) error {
	return nil
}

func (p) ExecuteHookCleanUp(hook plugin.ExecutedHook) error {
	return nil
}

func main() {
	pluginMap := map[string]hplugin.Plugin{
		"airdrop": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
