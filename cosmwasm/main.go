package main

import (
	"embed"
	"encoding/gob"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
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

	err = modifyFiles(p.chainName)
	if err != nil {
		return err
	}
	/*
		err = modifyAppGo(p.chainName)
		if err != nil {
			return err
		}
	*/

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

func modifyFiles(chainName string) error {
	files := []struct {
		template, outDir, outFile string
	}{

		{
			template: "app_47_3.go.plush",
			outDir:   "app",
			outFile:  "app.go",
		},
	}

	for _, f := range files {

		err := modifyFilesHelper(f.outDir, f.outFile, chainName)
		if err != nil {
			return err
		}
	}

	return nil
}

// Install dependencies
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

func modifyAppGoHelper(inputFilename, outputDir, outputFilename string, chainName string) error {

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
	//ctx.Set("ModulePath", chainName)
	//ctx.Set("BinaryNamePrefix", chainName)
	//ctx.Set("AddressPrefix", "cosmos")

	// Render the Plush template using the sourceContent as the template string
	renderedContent, err := plush.Render(string(sourceContent), ctx)
	if err != nil {
		return err
	}
	renderedContent = html.UnescapeString(renderedContent)

	// Create a new file with the rendered content
	g := genny.New()
	g.File(genny.NewFileS(outputPath, renderedContent))

	// Write the embedded content to the output file
	err = os.WriteFile(outputPath, []byte(renderedContent), 0o644)
	if err != nil {
		return err
	}

	fmt.Println("Created", filepath.Join(outputDir, outputFilename))

	// Parse the Go file to get the AST.
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filepath.Join(outputDir, outputFilename), nil, parser.ParseComments)
	if err != nil {
		return err
	}

	const maxFiles = 23
	placeholderContents := make([][]byte, maxFiles)

	for i := 1; i <= maxFiles; i++ {
		content, err := templates.ReadFile(fmt.Sprintf("placeholder_code/app%d.plush", i))
		if err != nil {
			return err
		}
		placeholderContents[i-1] = content
		fmt.Printf("Content of placeholder %d: %s\n", i, string(content))
	}

	// Traverse the AST and identify the places where you want to modify the code.
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			// Parse the chunk string into its own AST.
			chunkAST, err := parser.ParseFile(fset, "", placeholderContents[2], parser.ParseComments)
			if err != nil {
				return false
			}
			// Extract the declarations from the chunk's AST.
			chunkDecls := chunkAST.Decls
			for i, decl := range x.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == "getGovProposalHandlers" {
					// Insert the chunkDecls before the current function declaration.
					x.Decls = append(x.Decls[:i], append(chunkDecls, x.Decls[i:]...)...)
					break
				}
			}

		case *ast.GenDecl:
			// For import block
			if x.Tok == token.IMPORT {
				// Insert placeholderContents[0] at the beginning of the import block
				spec := &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: string(placeholderContents[0]),
					},
				}
				// Insert the new spec at the beginning of the Specs slice
				newSpecs := []ast.Spec{spec}
				for _, s := range x.Specs {
					newSpecs = append(newSpecs, s)
				}
				x.Specs = newSpecs

				// Append placeholderContent2 at the end of the import block
				spec2 := &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: string(placeholderContents[1]),
					},
				}
				x.Specs = append(x.Specs, spec2)
			}

		}
		return true
	})

	// Construct the output file path.
	outputPath = filepath.Join(outputDir, outputFilename)

	// Write the modified AST back to the file.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = format.Node(outputFile, fset, node)
	if err != nil {
		return err
	}

	fmt.Println("Modified", outputPath)

	return nil

}

func replaceInString(s, old, new string) string {
	return string([]byte(s))
}

// Helper to find index of decl in slice
func indexOfDecl(decls []ast.Decl, decl ast.Decl) int {
	for i, d := range decls {
		if d == decl {
			return i
		}
	}
	return -1
}

func modifyAppGo(chainName string) error {
	files := []struct {
		template, outDir, outFile string
	}{

		{
			template: "app_47_3.go.plush",
			outDir:   "app",
			outFile:  "app.go",
		},
	}

	for _, f := range files {

		err := modifyAppGoHelper(f.template, f.outDir, f.outFile, chainName)
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
