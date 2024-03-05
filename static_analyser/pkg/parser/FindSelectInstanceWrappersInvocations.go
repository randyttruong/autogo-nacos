package parser

import (
	"fmt"
	"go/ast"
	t "static_analyser/pkg/types"
	"strings"
)

func FindSelectInstanceWrappersInvocations(node ast.Node, wrapper t.ServiceDiscoveryWrapper, service string) []string {
	// FindSelectInstanceWrappersInvocations is a function that finds the invocation of the wrappers for service discovery and resolves the arguments for serviceName.
	//
	// node: The root node of the AST.
	// wrapper: The ServiceDiscoveryWrapper struct that contains the wrapper function and the arguments to resolve.
	// service: The name of the service.
	//
	// Returns:
	// A slice of service names.

	handleBasicLit := func(arg ast.Expr) string {
		// handleBasicLit is a closure that processes an *ast.BasicLit node and returns its value as a string.
		//
		// arg: The *ast.BasicLit node to process.
		//
		// Returns:
		// The value of the *ast.BasicLit node as a string if the node is of type *ast.BasicLit. If the node is not of type *ast.BasicLit, it returns "nil".

		// Check if the argument is a basic literal
		if lit, ok := arg.(*ast.BasicLit); ok {
			return lit.Value
		}
		return "nil"
	}

	resolveArgument := func(arg interface{}, args []string) string {
		// resolveArgument is a closure that processes an argument and returns its value as a string.
		//
		// arg: The argument to process. It can be of type string or t.WrapperParams.
		// args: A slice of string arguments from the wrapper function.
		//
		// Returns:
		// The value of the argument as a string. If the argument is of type string, it returns the argument itself. If the argument is of type t.WrapperParams, it returns the argument at the position specified in the t.WrapperParams struct from the args slice. If the argument is of neither type, it returns an empty string.
		//

		switch t := arg.(type) {
		// Check if the argument is a string or a t.WrapperParams
		case string:
			// If the argument is a string, return the argument itself
			return t
		case t.WrapperParams:
			// If the argument is a t.WrapperParams, return the argument at the position specified in the t.WrapperParams struct from the args slice
			return strings.ReplaceAll(args[t.Position], "\"", "")
		}
		// If the argument is of neither type, return an empty string
		return ""
	}

	wrapperName := wrapper.Wrapper
	serviceNames := []string{}

	// Inspect the AST for function calls
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		// Check if the node is a *ast.CallExpr
		case *ast.CallExpr:
			// Check if the function is the wrapper function
			if fun, ok := n.Fun.(*ast.Ident); ok {
				var args []string
				if fun.Name == wrapperName {
					// Iterate over the arguments of the function call
					for _, arg := range n.Args {
						args = append(args, handleBasicLit(arg))
					}
					fmt.Printf("%v", args)

					serviceName := resolveArgument(wrapper.ServiceName, args)

					serviceNames = append(serviceNames, serviceName)
				}
			}
		}
		return true
	})

	return serviceNames
}
