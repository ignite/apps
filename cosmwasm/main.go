package main

import (
	"embed"
	"encoding/gob"
	"fmt"
	"html"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/services/plugin"
)

//go:embed wasm-wiring/*
//go:embed placeholder_code/*
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

func (p) Manifest() (plugin.Manifest, error) {
	return plugin.Manifest{
		Name: "cosmwasm",
		// Add commands here
		Commands: []plugin.Command{
			// Example of a command
			{
				Use:               "cosmwasm",
				Short:             "The cosmwasm command is used for adding cosmwasm support for apps scaffolded with Ignite CLI",
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
	}, nil
}

func (p *p) Execute(cmd plugin.ExecutedCommand) error {

	// According to the number of declared commands, you may need a switch:
	switch cmd.Use {
	case "add":
		return p.handleAddCommand()
	default:
		return fmt.Errorf("unknown command: %s", cmd.Use)
	}
}

func (p *p) handleAddCommand() error {
	fmt.Println("Updating app files...")

	err := installDependencies()
	if err != nil {
		return err
	}

	err = createNewFiles(p.chainName)
	if err != nil {
		return err
	}

	return nil
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

func createFile(inputFilename, outputDir, outputFilename string, chainName string) error {

	// Load the embedded template files
	templatesFS, err := fs.Sub(templates, "wasm-wiring")
	if err != nil {
		return err
	}

	// Open the embedded input file
	inputFile, err := templatesFS.Open(inputFilename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Read content from the embedded file
	sourceContent, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	// Construct the output file path
	outputPath := filepath.Join(outputDir, outputFilename)

	// Create a Plush context and set variables
	ctx := plush.NewContext()
	ctx.Set("planet", chainName)

	// Read the content of a placeholder code file
	if inputFilename == "app.go.plush" {
		for i := 1; i <= 23; i++ {
			placeholderContent, err := templates.ReadFile(fmt.Sprintf("placeholder_code/app%d.plush", i))
			if err != nil {
				return err
			}
			ctx.Set(fmt.Sprintf("app%d", i), string(placeholderContent))
		}
	}

	// Render the Plush template using the sourceContent as the template string
	renderedContent, err := plush.Render(string(sourceContent), ctx)
	if err != nil {
		return err
	}
	renderedContent = html.UnescapeString(renderedContent) // this will do the job

	// Create a new file with the rendered content
	g := genny.New()
	g.File(genny.NewFileS(outputPath, renderedContent))

	// Write the embedded content to the output file
	err = os.WriteFile(outputPath, []byte(renderedContent), 0o644)
	if err != nil {
		return err
	}

	fmt.Println("Created", filepath.Join(outputDir, outputFilename))

	return nil
}

func createNewFiles(chainName string) error {
	files := []struct {
		template, outDir, outFile string
	}{
		{
			template: "ante.go.plush",
			outDir:   "app",
			outFile:  "ante.go",
		},
		{
			template: "app.go.plush",
			outDir:   "app",
			outFile:  "app.go",
		},
		{
			template: "wasm.go.plush",
			outDir:   "app",
			outFile:  "wasm.go",
		},
		{
			template: "simulation_test.go.plush",
			outDir:   "app",
			outFile:  "simulation_test.go",
		},
		{
			template: "network.go.plush",
			outDir:   "testutil/network",
			outFile:  "network.go",
		},
		{
			template: "root.go.plush",
			outDir:   filepath.Join("cmd", chainName+"d", "cmd"),
			outFile:  "root.go",
		},
	}

	for _, f := range files {
		err := createFile(f.template, f.outDir, f.outFile, chainName)
		if err != nil {
			return err
		}
	}

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
