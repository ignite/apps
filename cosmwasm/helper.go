package main

import (
	"bytes"
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

	const maxFiles = 24
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
		case *ast.FuncDecl:
			// New consolidated code handling multiple checks within FuncDecl nodes
			if x.Name.Name == "New" || x.Name.Name == "main" { // Assuming main is also of interest
				newList := []ast.Stmt{}
				for _, stmt := range x.Body.List {
					remove := false

					// Check if it's an expression statement.
					if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
						if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
							if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								// Check if the function call is "SetLegacyRouter"
								if ident, ok := selExpr.X.(*ast.Ident); ok {
									if ident.Name == "govKeeper" && selExpr.Sel.Name == "SetLegacyRouter" {
										remove = true // Remove govKeeper.SetLegacyRouter
									}
									if ident.Name == "app" && selExpr.Sel.Name == "SetAnteHandler" {
										remove = true // Remove app.SetAnteHandler
									}
								}
							}
						}
					}

					// Check for assignment statements, specifically the anteHandler creation
					if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
						if len(assignStmt.Rhs) == 1 {
							if callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
								if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
									if ident, ok := selExpr.X.(*ast.Ident); ok {
										if ident.Name == "ante" && selExpr.Sel.Name == "NewAnteHandler" {
											remove = true // Remove anteHandler assignment
										}
									}
								}
							}
						}
					}

					if !remove {
						newList = append(newList, stmt)
					}
				}
				x.Body.List = newList
			}
		}
		return true // Continue traversing the AST
	})

	// Traverse the AST to find the New() function.
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true // Not a FuncDecl, skip to next node
		}

		// Check if the function name is "New"
		if funcDecl.Name.Name != "New" {
			return true // Not the "New" function, skip to next node
		}

		// We are in the right function, now look for the specific line to remove.
		newList := []ast.Stmt{}
		for _, stmt := range funcDecl.Body.List {
			remove := false

			// Check if it's an expression statement.
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						// Check if the function call is "SetLegacyRouter"
						if ident, ok := selExpr.X.(*ast.Ident); ok {
							if ident.Name == "govKeeper" && selExpr.Sel.Name == "SetLegacyRouter" {
								remove = true // This is the line we want to remove.
							}
						}
					}
				}
			}

			if !remove {
				// If the statement should not be removed, include it in the new list.
				newList = append(newList, stmt)
			}
		}

		// Replace the function's statements with the new list.
		funcDecl.Body.List = newList
		return false // We found what we were looking for; no need to continue.
	})

	// From this ast.Inspect we inject logic
	ast.Inspect(node, func(n ast.Node) bool {
		//fmt.Printf("Node type: %T at position %v\n", n, fset.Position(n.Pos()))

		switch x := n.(type) {

		case *ast.GenDecl:
			// For import block
			if x.Tok == token.IMPORT {
				fmt.Println("Inside the import block.")

				// Insert app1.plush at the beginning of the import block
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

				// Append  app2.plush  at the end of the import block
				var insertAfterIndex int = -1
				for i, spec := range x.Specs {
					if importSpec, ok := spec.(*ast.ImportSpec); ok {
						// Check if the import path matches the one we want to insert after
						if importSpec.Path.Value == fmt.Sprintf("\"%s/x/%s/types\"", chainName, chainName) {
							insertAfterIndex = i
							break
						}
					}
				}

				// If the specific import is found, insert the new import after it
				if insertAfterIndex != -1 {
					specToInsert := &ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: string(placeholderContents[1]),
						},
					}
					// Insert the new spec after the found import
					x.Specs = append(x.Specs[:insertAfterIndex+1], append([]ast.Spec{specToInsert}, x.Specs[insertAfterIndex+1:]...)...)
				}
			}

		case *ast.File:
			//injecting app3.plush
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

		//injecting app4.plush
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

	//appending app6.plush and app7.plush
	ast.Inspect(node, func(n ast.Node) bool {

		fs := token.NewFileSet()

		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Name.Name != "App" {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Identify the correct position to insert fields
		position := -1
		for idx, field := range structType.Fields.List {
			if len(field.Names) > 0 && field.Names[0].Name == "ConsensusParamsKeeper" {
				position = idx
				break
			}
		}

		if position == -1 {
			fmt.Println("Failed to find the correct insertion position.")
			return false
		}

		// Parse the placeholder content into AST
		placeholderAST, err := parser.ParseFile(fs, "", "package main; type _dummyStruct struct {"+string(placeholderContents[5])+"}", 0)
		if err != nil {
			panic(err)
		}

		// Extract the fields from the parsed placeholder
		placeholderFields := placeholderAST.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List

		// Insert fields at the correct position
		structType.Fields.List = append(structType.Fields.List[:position+1], append(placeholderFields, structType.Fields.List[position+1:]...)...)

		// app7.plush
		// Parse the placeholder content into AST
		placeholderAST, err = parser.ParseFile(fs, "", "package main; type _dummyStruct struct {"+string(placeholderContents[6])+"}", 0)
		if err != nil {
			panic(err)
		}

		// Extract the fields from the parsed placeholder
		placeholderFields = placeholderAST.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List

		for idx, field := range structType.Fields.List {
			if len(field.Names) > 0 && field.Names[0].Name == "ScopedICAHostKeeper" {
				structType.Fields.List = append(structType.Fields.List[:idx+1], append(placeholderFields, structType.Fields.List[idx+1:]...)...)
				return false
			}
		}

		return false
	})
	/*
		// appending app8.plush andd app9.plush
		ast.Inspect(node, func(n ast.Node) bool {
			fs := token.NewFileSet()

			// Parse the placeholder content into AST
			placeholderAST, err := parser.ParseFile(fs, "", "package main; func _dummyFunc("+string(placeholderContents[7])+") {}", 0)
			if err != nil {
				panic(err)
			}

			// Check if the current node is a function declaration
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != "New" {
				return true
			}

				// Find the parameter 'loadLatest bool' and insert after it
				for i, param := range funcDecl.Type.Params.List {
					if len(param.Names) > 0 && param.Names[0].Name == "loadLatest" {
						// Insert new parameter after it
						funcDecl.Type.Params.List = append(funcDecl.Type.Params.List, nil)
						copy(funcDecl.Type.Params.List[i+2:], funcDecl.Type.Params.List[i+1:])
						funcDecl.Type.Params.List[i+1] = placeholderAST.Decls[0].(*ast.FuncDecl).Type.Params.List[0]
						break
					}

				}


			placeholderAST, err = parser.ParseFile(fs, "", "package main; func _dummyFunc("+string(placeholderContents[8])+") {}", 0)
			if err != nil {
				panic(err)
			}
			// Find the parameter 'appOpts servertypes.AppOptions,' and insert after it
			for i, param := range funcDecl.Type.Params.List {
				if len(param.Names) > 0 && param.Names[0].Name == "appOpts" {
					// Insert new parameter after it
					funcDecl.Type.Params.List = append(funcDecl.Type.Params.List, nil)
					copy(funcDecl.Type.Params.List[i+2:], funcDecl.Type.Params.List[i+1:])
					funcDecl.Type.Params.List[i+1] = placeholderAST.Decls[0].(*ast.FuncDecl).Type.Params.List[0]
					break
				}
			}

			return true
		})

	*/
	// appending app10.plush ()
	ast.Inspect(node, func(n ast.Node) bool {
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		if len(assignStmt.Rhs) == 0 {
			return true
		}

		callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}

		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if selExpr.Sel.Name == "NewKVStoreKeys" {
				fmt.Println("Found NewKVStoreKeys")

				for idx, arg := range callExpr.Args {

					// Convert the argument to its string representation
					var buf bytes.Buffer
					err := format.Node(&buf, fset, arg)
					if err != nil {
						panic(err)
					}

					if buf.String() == chainName+"moduletypes.StoreKey" {

						// Create a new AST node for our placeholder content.
						newKey := &ast.Ident{
							Name: string(placeholderContents[9]),
						}
						// Insert our new key right after the located key.
						callExpr.Args = append(callExpr.Args[:idx+1], append([]ast.Expr{newKey}, callExpr.Args[idx+1:]...)...)
						return false
					}
				}
			}
		}

		return true
	})

	// appending app12.plush (+ content from app11.plush)
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			// Debugging: Print number of declarations found.
			fmt.Printf("Found %d declarations in the AST.\n", len(x.Decls))

			// Parse the chunk string into its own AST as statements.
			chunkAST, err := parser.ParseFile(fset, "", "package main\n func _(){"+string(placeholderContents[11])+"}", parser.ParseComments)
			if err != nil {
				fmt.Printf("Parse error: %v\n", err)
				return false
			}

			// Debugging: Print the parsed statements.
			fmt.Printf("Parsed the placeholder contents into %d declarations.\n", len(chunkAST.Decls))

			// We're looking for the first function within the parsed AST since we used a dummy function to parse statements.
			var chunkStmts []ast.Stmt
			for _, decl := range chunkAST.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok {
					chunkStmts = funcDecl.Body.List // Extract the statements from the dummy function's body.
					break
				}
			}

			// Debugging: Print the number of statements extracted.
			fmt.Printf("Extracted %d statements from the placeholder contents.\n", len(chunkStmts))

			// Now we traverse the actual AST of our target file.
			for i, decl := range x.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if !ok || funcDecl.Name.Name != "New" {
					continue
				}
				fmt.Println("Found 'New' function declaration")

				for stmtIdx, stmt := range funcDecl.Body.List {
					assignStmt, ok := stmt.(*ast.AssignStmt)
					if !ok {
						continue
					}

					// Check if this assignment statement is the one we're looking for.
					if len(assignStmt.Lhs) == 1 {
						if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok && ident.Name == "icaHostIBCModule" {
							fmt.Printf("Found the 'icaHostIBCModule' assignment at statement index %d\n", stmtIdx)

							// Insert our chunk statements after the identified statement.
							funcDecl.Body.List = append(funcDecl.Body.List[:stmtIdx+1], append(chunkStmts, funcDecl.Body.List[stmtIdx+1:]...)...)

							// Debugging: Indicate insertion point.
							fmt.Printf("Inserted chunk statements after 'icaHostIBCModule' assignment.\n")

							// Update the AST in the original file node.
							x.Decls[i] = funcDecl
							return false // Stop the inspection if insertion is done.
						}

					}

				}

			}

		}
		return true
	})

	// appending app13.plush
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			// Debugging: Print number of declarations found.
			fmt.Printf("Found %d declarations in the AST.\n", len(x.Decls))

			// Parse the chunk string into its own AST as statements.
			chunkAST, err := parser.ParseFile(fset, "", "package main\n func _(){"+string(placeholderContents[12])+"}", parser.ParseComments)
			if err != nil {
				fmt.Printf("Parse error: %v\n", err)
				return false
			}

			// Debugging: Print the parsed statements.
			fmt.Printf("Parsed the placeholder contents into %d declarations.\n", len(chunkAST.Decls))

			// We're looking for the first function within the parsed AST since we used a dummy function to parse statements.
			var chunkStmts []ast.Stmt
			for _, decl := range chunkAST.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok {
					chunkStmts = funcDecl.Body.List // Extract the statements from the dummy function's body.
					break
				}
			}

			// Now we traverse the actual AST of our target file.
			for i, decl := range x.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if !ok || funcDecl.Name.Name != "New" {
					continue
				}
				fmt.Println("Found 'New' function declaration")

				for stmtIdx, stmt := range funcDecl.Body.List {
					// Check if this statement is an expression statement.
					exprStmt, ok := stmt.(*ast.ExprStmt)
					if !ok {
						continue
					}

					// Check if the expression is a call expression.
					callExpr, ok := exprStmt.X.(*ast.CallExpr)
					if !ok {
						continue
					}

					// Now we need to determine if this call expression is `app.CapabilityKeeper.Seal()`
					selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
					if !ok {
						continue
					}

					// Check if the selector expression matches our target method call.
					if ident, ok := selectorExpr.X.(*ast.SelectorExpr); ok {
						if xIdent, ok := ident.X.(*ast.Ident); ok && xIdent.Name == "app" &&
							ident.Sel.Name == "CapabilityKeeper" && selectorExpr.Sel.Name == "Seal" {

							fmt.Printf("Found the 'app.CapabilityKeeper.Seal()' call at statement index %d\n", stmtIdx)

							// Insert our chunk statements after the identified statement.
							funcDecl.Body.List = append(funcDecl.Body.List[:stmtIdx+1], append(chunkStmts, funcDecl.Body.List[stmtIdx+1:]...)...)

							// Debugging: Indicate insertion point.
							fmt.Printf("Inserted chunk statements after 'app.CapabilityKeeper.Seal()'.\n")

							// Update the AST in the original file node.
							x.Decls[i] = funcDecl
							return false // Stop the inspection if insertion is done.
						}
					}
				}
			}
		}
		return true
	})

	// appending app14.plush
	// Find the router config assignment
	var ibcRouterDecl *ast.AssignStmt
	ast.Inspect(node, func(n ast.Node) bool {
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			// Check that we are assigning to a variable named "ibcRouter"
			if len(assignStmt.Lhs) > 0 {
				if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok && ident.Name == "ibcRouter" {
					ibcRouterDecl = assignStmt
					return false // Stop inspecting further
				}
			}
		}
		return true // Continue inspecting
	})

	// Insert code after ica route
	if ibcRouterDecl != nil {
		// Code to insert
		newCode := string(placeholderContents[13])

		// Parse the insertion as an expression
		newExpr, err := parser.ParseExpr(newCode)
		if err != nil {
			panic(err) // Handle the error appropriately
		}

		// Type assert the last argument of the router declaration to *ast.CallExpr
		lastCallExpr, ok := ibcRouterDecl.Rhs[0].(*ast.CallExpr)
		if !ok {
			panic("right-hand side of assignment is not a call expression")
		}

		// Create a new call expression appending the new route
		// Assuming that lastCallExpr is the end of the call chain where we want to add the new route
		newCallExpr := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   lastCallExpr,
				Sel: newExpr.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel,
			},
			Args: newExpr.(*ast.CallExpr).Args,
		}

		// Update the AST by setting the new call expression
		ibcRouterDecl.Rhs[0] = newCallExpr
	}

	//appending app15.plush
	// Define the new code to be injected
	newCode := string(placeholderContents[14])
	// Parse the new code to get an AST node
	newStmt, err := parser.ParseExpr(newCode)
	if err != nil {
		fmt.Printf("Could not parse new code: %v\n", err)
		return err
	}

	// Traverse the AST to find the module manager initialization

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for an assignment statement (which could be setting up the module manager)
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true // continue searching
		}

		// Check if the right-hand side of the assignment is a call to module.NewManager
		callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true // continue searching
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || selExpr.Sel.Name != "NewManager" {
			return true // continue searching
		}

		// Now we have the call to module.NewManager, let's find ibc.NewAppModule
		for i, arg := range callExpr.Args {
			call, ok := arg.(*ast.CallExpr)
			if !ok {
				continue // not a call expression, skip
			}

			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "ibc" && sel.Sel.Name == "NewAppModule" {
					// Found ibc.NewAppModule, insert new module before it
					callExpr.Args = append(callExpr.Args[:i], append([]ast.Expr{newStmt}, callExpr.Args[i:]...)...)

					return false // stop searching
				}
			}
		}

		return true
	})

	//appending app16-17 plush
	// Define the new code to be injected
	newCode = string(placeholderContents[15])

	// Parse the new code to get an AST node
	newExpr, err := parser.ParseExpr(newCode)
	if err != nil {
		fmt.Printf("Could not parse new code: %v\n", err)
		return err
	}

	// Traverse the AST to find the SetOrderBeginBlockers call

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for a call expression
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true // continue searching
		}

		// Check if the function called is SetOrderBeginBlockers
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if selExpr.Sel.Name == "SetOrderBeginBlockers" {
				// Append the new module name to the end of the call's arguments
				callExpr.Args = append(callExpr.Args, newExpr)

				return false // stop searching
			}
		}
		// Check if the function called is SetOrderBeginBlockers
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if selExpr.Sel.Name == "SetOrderEndBlockers" {
				// Append the new module name to the end of the call's arguments
				callExpr.Args = append(callExpr.Args, newExpr)

				return false // stop searching
			}
		}

		return true
	})

	//appending app18 plush
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for Function Declarations
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check if the function name is New
		if funcDecl.Name.Name != "New" {
			return true
		}

		// Now we are inside the New function, look for the genesisModuleOrder variable
		for _, stmt := range funcDecl.Body.List {
			// Look for Assign Statements
			assignStmt, ok := stmt.(*ast.AssignStmt)
			if !ok {
				continue
			}

			// Check if the LHS of the assignment has the identifier we are looking for
			for _, lhs := range assignStmt.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || ident.Name != "genesisModuleOrder" {
					continue
				}

				// We found the genesisModuleOrder assignment, now add the new element
				if compLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit); ok {
					compLit.Elts = append(compLit.Elts, newExpr)

					return false // stop searching
				}
			}
		}

		return true
	})

	//injecting app19.plush
	// Define the new code to be injected
	newCode = string(placeholderContents[18])

	// Parse the new code to get an AST node
	newExpr, err = parser.ParseExpr(newCode)
	if err != nil {
		fmt.Printf("Could not parse new code: %v\n", err)
		return err
	}

	// Convert the expression to a statement
	newStmt1 := &ast.ExprStmt{X: newExpr}

	// Traverse the AST to find the New function and then the SetInitChainer call

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for Function Declarations
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check if the function name is New
		if funcDecl.Name.Name != "New" {
			return true
		}

		// Now we are inside the New function, look for the SetInitChainer call
		for i, stmt := range funcDecl.Body.List {
			// Look for Expression Statements
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}

			// Check if the expression is a call to SetInitChainer
			callExpr, ok := exprStmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}

			// Check if the function being called is SetInitChainer
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "app" && selExpr.Sel.Name == "SetInitChainer" {
					// We found the SetInitChainer call, now inject the new statement before it
					funcDecl.Body.List = append(funcDecl.Body.List[:i], append([]ast.Stmt{newStmt1}, funcDecl.Body.List[i:]...)...)

					return false // stop searching
				}
			}
		}

		return true
	})

	//injecting app20.plush
	// Parse the new code to get an AST node
	newStmts, err := parser.ParseFile(fset, "", "package main\nfunc _() {"+string(placeholderContents[19])+"}", parser.ParseComments)
	if err != nil {
		fmt.Printf("Could not parse new code: %v\n", err)
		return err
	}

	// Extract the block statement from the parsed code
	newBlockStmt := newStmts.Decls[0].(*ast.FuncDecl).Body

	// Traverse the AST to find the New function and then the SetEndBlocker call
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for Function Declarations
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check if the function name is New
		if funcDecl.Name.Name != "New" {
			return true
		}

		// Now we are inside the New function, look for the SetEndBlocker call
		for i, stmt := range funcDecl.Body.List {
			// Look for Expression Statements
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}

			// Check if the expression is a call to SetEndBlocker
			callExpr, ok := exprStmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}

			// Check if the function being called is SetEndBlocker
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "app" && selExpr.Sel.Name == "SetEndBlocker" {
					// We found the SetEndBlocker call, now inject the new block statement after it
					funcDecl.Body.List = append(funcDecl.Body.List[:i+1], append(newBlockStmt.List, funcDecl.Body.List[i+1:]...)...)

					return false // stop searching
				}
			}
		}

		return true
	})

	//injecting app21.plush

	// Parse the new code to get an AST node
	newStmts, err = parser.ParseFile(fset, "", "package main\nfunc _() {"+string(placeholderContents[20])+"}", parser.ParseComments)
	if err != nil {
		fmt.Printf("Could not parse new code: %v\n", err)
		return err
	}

	// Extract the block statement from the parsed code
	newBlockStmt = newStmts.Decls[0].(*ast.FuncDecl).Body

	// Traverse the AST to find the if loadLatest statement
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for If Statements
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		// Check if the condition is a comparison with loadLatest
		if ident, ok := ifStmt.Cond.(*ast.Ident); ok && ident.Name == "loadLatest" {
			// We found the if loadLatest statement, now inject the new block statement at the end of it
			ifStmt.Body.List = append(ifStmt.Body.List, newBlockStmt.List...)
			return false // stop searching
		}

		return true
	})

	//injecting app22.plush
	// Parse the code chunk as a statement by wrapping it in a function
	wrappedCodeChunk := fmt.Sprintf("package main\nfunc _() {\n%s\n}", string(placeholderContents[21]))
	tempFile, err := parser.ParseFile(fset, "", wrappedCodeChunk, parser.ParseComments)
	if err != nil {
		fmt.Printf("Could not parse wrapped code chunk: %v\n", err)
		return err
	}
	injectStmt := tempFile.Decls[0].(*ast.FuncDecl).Body.List[0]

	// Traverse the AST to find the New function
	ast.Inspect(node, func(n ast.Node) bool {
		// Check if this is a function declaration
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "New" {
			return true // not the "New" function, skip to the next node
		}

		// Look for the specific assignment statement
		for i, stmt := range funcDecl.Body.List {
			if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
				if selectorExpr, ok := assignStmt.Lhs[0].(*ast.SelectorExpr); ok {
					if ident, ok := selectorExpr.X.(*ast.Ident); ok && ident.Name == "app" && selectorExpr.Sel.Name == "ScopedTransferKeeper" {
						// Found the assignment, now inject the new statement after it
						funcDecl.Body.List = append(funcDecl.Body.List[:i+1], append([]ast.Stmt{injectStmt}, funcDecl.Body.List[i+1:]...)...)
						return false // we've done our injection, no need to traverse further
					}
				}
			}
		}
		return true
	})

	//injecting app23.plush
	// Parse the function to inject as a declaration
	wrappedFunctionToInject := fmt.Sprintf("package main\n%s", string(placeholderContents[22]))
	tempFile, err = parser.ParseFile(fset, "", wrappedFunctionToInject, parser.ParseComments)
	if err != nil {
		fmt.Printf("Could not parse wrapped function to inject: %v\n", err)
		return err
	}
	injectFuncDecl := tempFile.Decls[0].(*ast.FuncDecl)

	// Find the end of the New function and inject the new function after it
	for i, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == "New" {
			// Inject the new function declaration after the New function
			node.Decls = append(node.Decls[:i+1], append([]ast.Decl{injectFuncDecl}, node.Decls[i+1:]...)...)
			break
		}
	}

	//injecting app24.plush
	// Parse the function to inject as a declaration
	wrappedFunctionToInject = fmt.Sprintf("package main\n%s", string(placeholderContents[23]))
	tempFile, err = parser.ParseFile(fset, "", wrappedFunctionToInject, parser.ParseComments)
	if err != nil {
		fmt.Printf("Could not parse wrapped function to inject: %v\n", err)
		return err
	}
	injectFuncDecl = tempFile.Decls[0].(*ast.FuncDecl)

	// Find the end of the New function and inject the new function after it
	for i, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == "New" {
			// Inject the new function declaration after the New function
			node.Decls = append(node.Decls[:i+1], append([]ast.Decl{injectFuncDecl}, node.Decls[i+1:]...)...)
			break
		}
	}

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

// helper for app4.plush injections
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
