package main

import (
	"encoding/gob"
	"fmt"
	"path/filepath"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/plugin"
)

func init() {
	gob.Register(plugin.Manifest{})
	gob.Register(plugin.ExecutedCommand{})
	gob.Register(plugin.ExecutedHook{})
}

type p struct{}

func (p) Manifest() (plugin.Manifest, error) {
	return plugin.Manifest{
		Name: "cosmwasm-plugin",
		// Add commands here
		Commands: []plugin.Command{
			// Example of a command
			{
				Use:   "cosmwasm-plugin",
				Short: "Explain what the command is doing...",
				Long:  "Long description goes here...",
				Flags: []plugin.Flag{
					{Name: "my-flag", Type: plugin.FlagTypeString, Usage: "my flag description"},
				},
				PlaceCommandUnder: "ignite",
				// Examples of adding subcommands:
				/*
					Commands: []plugin.Command{
						{Use: "add"},
						{Use: "list"},
						{Use: "delete"},
					},
				*/
			},
		},
		// Add hooks here
		Hooks: []plugin.Hook{},
		SharedHost: false,
	}, nil
}

func (p) Execute(cmd plugin.ExecutedCommand) error {
	// TODO: write command execution here
	fmt.Printf("Hello I'm the cosmwasm-plugin plugin\n")
	fmt.Printf("My executed command: %q\n", cmd.Path)
	fmt.Printf("My args: %v\n", cmd.Args)
	myFlag, _ := cmd.Flags().GetString("my-flag")
	fmt.Printf("My flags: my-flag=%q\n", myFlag)
	fmt.Printf("My config parameters: %v\n", cmd.With)

	// This is how the plugin can access the chain:
	// c, err := getChain(cmd)

	// According to the number of declared commands, you may need a switch:
	/*
		switch cmd.Use {
		case "add":
			fmt.Println("Adding stuff...")
		case "list":
			fmt.Println("Listing stuff...")
		case "delete":
			fmt.Println("Deleting stuff...")
		}
	*/
	return nil
}

func (p) ExecuteHookPre(hook plugin.ExecutedHook) error {
	fmt.Printf("Executing hook pre %q\n", hook.Name)
	return nil
}

func (p) ExecuteHookPost(hook plugin.ExecutedHook) error {
	fmt.Printf("Executing hook post %q\n", hook.Name)
	return nil
}

func (p) ExecuteHookCleanUp(hook plugin.ExecutedHook) error {
	fmt.Printf("Executing hook cleanup %q\n", hook.Name)
	return nil
}

func getChain(cmd plugin.ExecutedCommand, chainOption ...chain.Option) (*chain.Chain, error) {
	var (
		home, _ = cmd.Flags().GetString("home")
		path, _ = cmd.Flags().GetString("path")
	)
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return chain.New(absPath, chainOption...)
}

func main() {
	pluginMap := map[string]hplugin.Plugin{
		"cosmwasm-plugin": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
