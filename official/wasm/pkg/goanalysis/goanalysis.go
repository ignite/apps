package goanalysis

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

// AppendImports inserts import statements into Go source code content.
func AppendImports(fileContent string, importStatements ...string) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Check if the import already exists.
	existImports := make(map[string]struct{})
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT || len(genDecl.Specs) == 0 {
			continue
		}

		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			existImports[importSpec.Path.Value] = struct{}{}
		}
	}

	newSpecs := make([]ast.Spec, 0)
	for _, importStatement := range importStatements {
		// Check if the import already exists.
		if _, ok := existImports[`"`+importStatement+`"`]; ok {
			continue
		}
		// Create a new import spec.
		newSpecs = append(newSpecs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"` + importStatement + `"`,
			},
		})
	}

	if len(newSpecs) == 0 {
		// No new imports to add.
		return fileContent, nil
	}

	// Create a new import declaration.
	newImportDecl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: newSpecs,
	}

	// Insert the new import declaration at the beginning of the file.
	newDecls := append([]ast.Decl{newImportDecl}, f.Decls...)

	// Update the file's declarations.
	f.Decls = newDecls

	// Write the modified AST to a buffer.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// AppendCode inserts code before the end or return of a function in Go source code content.
func AppendCode(fileContent, functionName, codeToInsert string) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Check if the function has the name you want to replace.
			if funcDecl.Name.Name == functionName {
				// Insert the code before the end or return of the function.
				insertionExpr, err := parser.ParseExpr(codeToInsert)
				if err != nil {
					return false
				}
				// Check if there is a return statement in the function.
				if len(funcDecl.Body.List) > 0 {
					lastStmt := funcDecl.Body.List[len(funcDecl.Body.List)-1]
					switch lastStmt.(type) {
					case *ast.ReturnStmt:
						// If there is a return, insert before it.
						funcDecl.Body.List = append(funcDecl.Body.List[:len(funcDecl.Body.List)-1], &ast.ExprStmt{X: insertionExpr}, lastStmt)
					default:
						// If there is no return, insert at the end of the function body.
						funcDecl.Body.List = append(funcDecl.Body.List, &ast.ExprStmt{X: insertionExpr})
					}
				} else {
					// If there are no statements in the function body, insert at the end of the function body.
					funcDecl.Body.List = append(funcDecl.Body.List, &ast.ExprStmt{X: insertionExpr})
				}
				found = true
				return false
			}
		}
		return true
	})

	if !found {
		return "", errors.Errorf("function %s not found", functionName)
	}

	// Write the modified AST to a buffer.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ReplaceReturn replaces return statements in a Go function with a new return statement.
func ReplaceReturn(fileContent, functionName, newReturnStatement string) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Check if the function has the name you want to replace.
			if funcDecl.Name.Name == functionName {
				// Replace the return statements.
				for _, stmt := range funcDecl.Body.List {
					if retStmt, ok := stmt.(*ast.ReturnStmt); ok {
						// Remove existing return statements.
						retStmt.Results = nil

						// Parse the new return statement.
						var buf bytes.Buffer
						buf.WriteString(newReturnStatement)
						returnExpr, err := parser.ParseExpr(buf.String())
						if err != nil {
							return false
						}
						// Add the new return statement.
						retStmt.Results = []ast.Expr{returnExpr}

						//retStmt.Results = append(retStmt.Results, newRetExpr)
					}
				}
				found = true
				return false
			}
		}
		return true
	})

	if !found {
		return "", errors.Errorf("function %s not found", functionName)
	}

	// Write the modified AST to a buffer.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}
