package main

import (
	"embed"
	"encoding/gob"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/services/plugin"
)

//go:embed wasm-wiring/*
var templates embed.FS // Embedded template files

func init() {
	gob.Register(plugin.Manifest{})
	gob.Register(plugin.ExecutedCommand{})
	gob.Register(plugin.ExecutedHook{})
}

type p struct {
	chainName string
}

// NewPlugin creates a new plugin instance with the given chainName and chainNamed values.
func NewPlugin(chainName string) *p {
	return &p{
		chainName: chainName,
	}
}

// ChainName returns the chainName of the plugin.
func (p *p) ChainName() string {
	return p.chainName
}

// ChainNamed derives and returns the chainNamed of the plugin.
func (p *p) ChainNamed() string {
	return p.chainName + "d"
}
func (p) Manifest() (plugin.Manifest, error) {
	return plugin.Manifest{
		Name: "cosmwasm",
		// Add commands here
		Commands: []plugin.Command{
			// Example of a command
			{
				Use:   "cosmwasm",
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
						PlaceCommandUnder: "ignite cosmwasm",
					},
				},
			},
		},
		// Add hooks here
		Hooks:      []plugin.Hook{},
		SharedHost: false,
	}, nil
}

func (p *p) Execute(cmd plugin.ExecutedCommand) error {

	fmt.Printf("Hello I'm the cosmwasm plugin\n")

	// According to the number of declared commands, you may need a switch:

	switch cmd.Use {
	case "add":
		return handleAddCommand(p)
	default:
		return fmt.Errorf("unknown command: %s", cmd.Use)
	}
}

func handleAddCommand(p *p) error {
	fmt.Println("Adding stuff...")

	err := installDependencies()
	if err != nil {
		return err
	}

	err = createNewFiles(p)
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

func createFile(inputFilename, outputDir, outputFilename string, p *p) error {

	// Load the embedded template files
	templatesFS, err := fs.Sub(templates, "wasm-wiring")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Open the embedded input file
	inputFile, err := templatesFS.Open(inputFilename)

	if err != nil {
		return err
	}
	defer inputFile.Close() // This line schedules the file to be closed when the function returns

	// Read content from the embedded file
	sourceContent, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	// Construct the output file path
	outputPath := filepath.Join(outputDir, outputFilename)

	// Write the embedded content to the output file
	err = ioutil.WriteFile(outputPath, sourceContent, 0644)
	if err != nil {
		return err
	}

	// Replace chainName in the output file
	err = replaceWordsInFile(outputPath, "planet", p.ChainName())
	if err != nil {
		return err
	}

	fmt.Printf("Created %s in %s\n", outputFilename, outputDir)
	return nil
}

func createNewFiles(p *p) error {
	files := []struct {
		template, outDir, outFile string
	}{
		{
			template: "ante.txt",
			outDir:   "app",
			outFile:  "ante.go",
		},
		{
			template: "wasm.txt",
			outDir:   "app",
			outFile:  "wasm.go",
		},
		{
			template: "app.txt",
			outDir:   "app",
			outFile:  "app.go",
		},
		{
			template: "simulation_test.txt",
			outDir:   "app",
			outFile:  "simulation_test.go",
		},
		{
			template: "network.txt",
			outDir:   "testutil/network",
			outFile:  "network.go",
		},
		{
			template: "root.txt",
			outDir:   filepath.Join("cmd", p.ChainNamed(), "cmd"),
			outFile:  "root.go",
		},
		// TODO: Add any other templates as needed
	}

	for _, f := range files {
		err := createFile(f.template, f.outDir, f.outFile, p)
		if err != nil {
			return err
		}
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
	chainName := filepath.Base(currentDir)

	// Handle potential cases where the chain name contains dashes or underscores
	chainName = strings.ReplaceAll(chainName, "-", "")
	chainName = strings.ReplaceAll(chainName, "_", "")

	pluginInstance := NewPlugin(chainName)

	pluginMap := map[string]hplugin.Plugin{
		"cosmwasm": &plugin.InterfacePlugin{Impl: pluginInstance},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
