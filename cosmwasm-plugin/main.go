package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/plugin"
)

var chainName string
var chainNamed string

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
				Commands: []plugin.Command{
					{
						Use:               "add",
						Short:             "Add a CosmWasm module to the chain",
						Long:              "This command will install dependencies, create new files, and modify existing files to add a CosmWasm module to the chain.",
						PlaceCommandUnder: "ignite cosmwasm-plugin",
					},
				},
			},
		},
		// Add hooks here
		Hooks:      []plugin.Hook{},
		SharedHost: false,
	}, nil
}

func (p) Execute(cmd plugin.ExecutedCommand) error {

	fmt.Printf("Hello I'm the cosmwasm-plugin\n")

	// This is how the plugin can access the chain:
	// c, err := getChain(cmd)

	// According to the number of declared commands, you may need a switch:

	switch cmd.Use {
	case "add":
		return handleAddCommand()
	default:
		return fmt.Errorf("unknown command: %s", cmd.Use)
	}

	return nil
}

func handleAddCommand() error {
	fmt.Println("Adding stuff...")

	err := installDependencies()
	if err != nil {
		return err
	}

	err = createNewFiles()
	if err != nil {
		return err
	}

	//err = modifyExistingFiles()
	//if err != nil {return err}

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

// For future use
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

func installDependencies() error {
	cmd := exec.Command("go", "get", "github.com/CosmWasm/wasmd/x/wasm@v0.41.0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to install dependency: %w", err)
	}
	return nil
}

func replaceWordsInFile(filePath, targetWord, replacement string) error {
	// Read content from the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Perform the replacement
	modifiedContent := strings.ReplaceAll(string(content), targetWord, replacement)

	// Write the modified content back to the file
	err = ioutil.WriteFile(filePath, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	//fmt.Printf("Replaced planet with %s\n", replacement)

	return nil
}

func createFile(inputFilename, outputDir, outputFilename string) error {
	sourcePath := filepath.Join("..", "cosmwasm-plugin", "wasm-wiring", inputFilename)

	sourceContent, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, outputFilename)

	err = ioutil.WriteFile(outputPath, sourceContent, 0644)
	if err != nil {
		return err
	}

	//replace chainName
	err = replaceWordsInFile(outputPath, "planet", chainName)
	if err != nil {
		return err
	}

	fmt.Printf("Created %s in %s\n", outputFilename, outputDir)
	return nil
}

func createNewFiles() error {
	// Create ante.go based on ante.txt content
	err := createFile("ante.txt", "app", "ante.go")
	if err != nil {
		return err
	}

	// Create ante2.go based on wasm.txt content
	err = createFile("wasm.txt", "app", "wasm.go")
	if err != nil {
		return err
	}

	// Create network.go based on network.txt content
	err = createFile("app.txt", "app", "app.go")
	if err != nil {
		return err
	}

	// Create network.go based on network.txt content
	err = createFile("simulation_test.txt", "app", "simulation_test.go")
	if err != nil {
		return err
	}

	// Create network.go based on network.txt content
	err = createFile("network.txt", "testutil/network", "network.go")
	if err != nil {
		return err
	}

	filePath := filepath.Join("cmd", chainNamed, "cmd")
	// Create network.go based on network.txt content
	err = createFile("root.txt", filePath, "root.go")
	if err != nil {
		return err
	}

	return nil
}

// For future versions..
func modifyExistingFiles() error {
	// Modify existingFile1.go
	content, err := os.ReadFile("existingFile1.go")
	if err != nil {
		return err
	}
	modifiedContent := string(content) + "\n// Additional content for existingFile1.go"
	err = os.WriteFile("existingFile1.go", []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	// Similarly, modify other existing files...

	return nil
}

func main() {

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Extract the chain name from the current directory
	chainName = filepath.Base(currentDir)

	// Handle potential cases where the chain name contains dashes or underscores
	chainName = strings.ReplaceAll(chainName, "-", "")
	chainName = strings.ReplaceAll(chainName, "_", "")
	chainNamed = chainName + "d"

	//fmt.Println("User's chain name:", chainName)

	pluginMap := map[string]hplugin.Plugin{
		"cosmwasm-plugin": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
