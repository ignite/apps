package goanalysis

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

// AppendImports appends import statements to the existing import block in Go source code content.
func AppendImports(fileContent string, importStatements ...string) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Find the existing import declaration.
	var importDecl *ast.GenDecl
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT || len(genDecl.Specs) == 0 {
			continue
		}
		importDecl = genDecl
		break
	}

	if importDecl == nil {
		// If no existing import declaration found, create a new one.
		importDecl = &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: make([]ast.Spec, 0),
		}
		f.Decls = append([]ast.Decl{importDecl}, f.Decls...)
	}

	// Check existing imports to avoid duplicates.
	existImports := make(map[string]struct{})
	for _, spec := range importDecl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}
		existImports[importSpec.Path.Value] = struct{}{}
	}

	// Add new import statements.
	for _, importStatement := range importStatements {
		var (
			impSplit   = strings.Split(importStatement, " ")
			importRepo = impSplit[len(impSplit)-1]
			importname = ""
		)
		if len(impSplit) > 1 {
			importname = impSplit[0]
		}

		// Check if the import already exists.
		if _, ok := existImports[strconv.Quote(importRepo)]; ok {
			continue
		}
		// Create a new import spec.
		spec := &ast.ImportSpec{
			Name: &ast.Ident{
				Name: importname,
			},
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(importRepo),
			},
		}
		importDecl.Specs = append(importDecl.Specs, spec)
	}

	// Format the modified AST.
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

	// Parse the Go code to insert.
	insertionExpr, err := parser.ParseExpr(codeToInsert)
	if err != nil {
		return "", err
	}

	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Check if the function has the name you want to replace.
			if funcDecl.Name.Name == functionName {
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
func ReplaceReturn(fileContent, functionName string, returnVars ...string) (string, error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	returnStmts := make([]ast.Expr, 0)
	for _, returnVar := range returnVars {
		// Parse the new return var to expression.
		newRetExpr, err := parser.ParseExpr(returnVar)
		if err != nil {
			return "", err
		}
		returnStmts = append(returnStmts, newRetExpr)
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
						// Add the new return statement.
						retStmt.Results = append(retStmt.Results, returnStmts...)
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

// ReplaceCode replaces a function implementation in Go source code content.
func ReplaceCode(fileContent, oldFunctionName, newFunction string) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Parse the content of the new function into an ast.File.
	newFuncContent := fmt.Sprintf("package p; func _() { %s }", strings.TrimSpace(newFunction))
	newFile, err := parser.ParseFile(fileSet, "", newFuncContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Check if the function has the name you want to replace.
			if funcDecl.Name.Name == oldFunctionName {
				// Take the body of the new function from the parsed file.
				newFunctionBody := newFile.Decls[0].(*ast.FuncDecl).Body
				// Replace the function body with the body of the new function.
				funcDecl.Body = newFunctionBody
				found = true
				return false
			}
		}
		return true
	})

	if !found {
		return "", fmt.Errorf("function %s not found in file content", oldFunctionName)
	}

	// Write the modified AST to a buffer.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}
