package main

import (
	"go/ast"
	"go/parser"
  "go/token"
	"strings"
)

// parseGoFile extracts function calls from a Go file.
// It returns a slice of function call names found in the file.
func parseGoFile(filePath string, function_names []string) ([]string, error) {
	var functionCalls []string

	// Prepare the file set
	fset := token.NewFileSet()
	// Parse the file
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
			return nil, err
	}

	// Walk through the AST
	ast.Inspect(node, func(n ast.Node) bool {
			// Check for function calls
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
					return true // continue walking the AST
			}
			switch fun := callExpr.Fun.(type) {
			case *ast.Ident: // Simple function calls
					functionCalls = append(functionCalls, fun.Name)
			case *ast.SelectorExpr: // Qualified (package) function calls, e.g., pkg.Func
					if id, ok := fun.X.(*ast.Ident); ok {
							fullName := id.Name + "." + fun.Sel.Name
							functionCalls = append(functionCalls, fullName)
					}
			}
			return true
	})

	return functionCalls, nil
}




// ArgumentType represents the type of an argument in a function call.
type ArgumentType string

const (
    Literal  ArgumentType = "Literal"
    Variable ArgumentType = "Variable"
    Unknown  ArgumentType = "Unknown"
)

// FunctionCallInfo contains information about a function call.
type FunctionCallInfo struct {
    FunctionName string
    Arguments    []ArgumentType
}

// parseGoFile analyzes function calls in a Go file, checking if their arguments are literals or variables.
func parseGoFile(filePath string, functionNames []string) ([]FunctionCallInfo, error) {
    var callsInfo []FunctionCallInfo
    funcMap := make(map[string]bool)
    for _, name := range functionNames {
        funcMap[name] = true
    }

    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
    if err != nil {
        return nil, err
    }

    ast.Inspect(node, func(n ast.Node) bool {
        callExpr, ok := n.(*ast.CallExpr)
        if !ok {
            return true
        }

        var funcName string
        switch fun := callExpr.Fun.(type) {
        case *ast.Ident:
            funcName = fun.Name
        case *ast.SelectorExpr:
            if id, ok := fun.X.(*ast.Ident); ok {
                funcName = id.Name + "." + fun.Sel.Name
            }
        }

        if funcMap[funcName] {
            var argTypes []ArgumentType
            for _, arg := range callExpr.Args {
                switch arg.(type) {
                case *ast.BasicLit: // Literal values (e.g., string, number)
                    argTypes = append(argTypes, Literal)
                case *ast.Ident: // Identifiers (e.g., variable names)
                    if strings.ToUpper(arg.(*ast.Ident).Name) == arg.(*ast.Ident).Name {
                        // Assuming constants are in uppercase, you might want to check for constants differently.
                        argTypes = append(argTypes, Literal)
                    } else {
                        argTypes = append(argTypes, Variable)
                    }
                default:
                    argTypes = append(argTypes, Unknown)
                }
            }
            callsInfo = append(callsInfo, FunctionCallInfo{
                FunctionName: funcName,
                Arguments:    argTypes,
            })
        }
        return true
    })

    return callsInfo, nil
}