package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func modifyFilesHelper(outputDir, outputFilename string, chainName string) error {
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
	}

	// First ast.Inspect for deleting code
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			// Remove existing variable named DefaultNodeHome
			if x.Tok == token.VAR {
				for i, spec := range x.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, name := range valueSpec.Names {
						if name.Name == "DefaultNodeHome" {
							// Remove this spec from the list
							x.Specs = append(x.Specs[:i], x.Specs[i+1:]...)
							return false
						}
					}
				}
			}
		}
		return true
	})

	// Second ast.Inspect for inserting new lines and any other logic
	ast.Inspect(node, func(n ast.Node) bool {
		//fmt.Printf("Node type: %T at position %v\n", n, fset.Position(n.Pos()))

		switch x := n.(type) {

		case *ast.GenDecl:
			// For import block
			if x.Tok == token.IMPORT {
				fmt.Println("Inside the import block.")

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

		}

		//append app4.plush
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || selExpr.Sel.Name != "NewBasicManager" {
			return true
		}

		modulesToInject, err := getModulesFromPlaceholder(fset, string(placeholderContents[3]))
		if err != nil {
			fmt.Println("Error parsing modules:", err)
			return false
		}

		fmt.Println("Modules to inject:", modulesToInject)
		// Inject the modules at the end
		callExpr.Args = append(callExpr.Args, modulesToInject...)

		fmt.Println("Injected the modules.")

		return true
	})

	//append app5.plush
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || valueSpec.Names[0].Name != "maccPerms" {
				continue
			}

			compLit, ok := valueSpec.Values[0].(*ast.CompositeLit)
			if !ok {
				return true
			}

			// Splitting the placeholder content
			parts := strings.SplitN(string(placeholderContents[4]), ": {", 2)
			if len(parts) != 2 {
				fmt.Println("Error: Placeholder content not in expected format")
				return false
			}

			moduleParts := strings.Split(parts[0], ".")
			if len(moduleParts) != 2 {
				fmt.Println("Error: Module part not in expected format")
				return false
			}

			permParts := strings.TrimRight(parts[1], "},")
			permItems := strings.Split(permParts, ",")

			// Construct the key and value
			key := &ast.SelectorExpr{
				X:   ast.NewIdent(moduleParts[0]),
				Sel: ast.NewIdent(moduleParts[1]),
			}

			var valueElts []ast.Expr
			for _, item := range permItems {
				perm := strings.Split(item, ".")
				if len(perm) != 2 {
					fmt.Println("Error: Permission not in expected format")
					return false
				}
				valueElts = append(valueElts, &ast.SelectorExpr{
					X:   ast.NewIdent(strings.TrimSpace(perm[0])),
					Sel: ast.NewIdent(strings.TrimSpace(perm[1])),
				})
			}

			value := &ast.CompositeLit{
				Elts: valueElts,
			}

			keyValueExpr := &ast.KeyValueExpr{
				Key:   key,
				Value: value,
			}

			// Append the new key-value pair to the end of the map literal
			compLit.Elts = append(compLit.Elts, keyValueExpr)

			fmt.Println("Injected into maccPerms.")
			return true
		}
		return true
	})

	// Write the modified AST back to the file.
	outputFile, err := os.Create(filepath.Join(outputDir, outputFilename))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = format.Node(outputFile, fset, node)
	if err != nil {
		return err
	}

	fmt.Println("Modified", filepath.Join(outputDir, outputFilename))

	return nil

}

func getModulesFromPlaceholder(fset *token.FileSet, content string) ([]ast.Expr, error) {
	moduleStrings := strings.Split(content, ",")
	var modules []ast.Expr

	for _, moduleStr := range moduleStrings {
		expr, err := parser.ParseExpr(strings.TrimSpace(moduleStr))
		if err != nil {
			fmt.Printf("Error parsing module '%s': %v\n", moduleStr, err)
			return nil, err
		}
		modules = append(modules, expr)
	}

	return modules, nil
}
