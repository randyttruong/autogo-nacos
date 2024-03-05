package parser

import (
	"go/ast"
	t "static_analyser/pkg/types"
	"strings"
)

// finds the invocation of the wrappers for register instance and resolves the arguments for serviceName, Ip, and Port
func FindRegisterInstanceWrapperInvocations(node ast.Node, wrapper t.RegisterInstanceWrapper, service string) ([]string, []t.ServiceInfo) {
	// FindRegisterInstanceWrapperInvocations finds the invocation of the wrappers for register instance and resolves the arguments for serviceName, Ip, and Port.
	//
	// node: The root node of the AST.
	// wrapper: The RegisterInstanceWrapper struct that contains the wrapper function and the arguments to resolve.
	// service: The name of the service.
	//
	// Returns:
	// A slice of service names and a slice of ServiceInfo structs. Each ServiceInfo struct contains the application name, IP, and port.
	//

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
	serviceInfos := []t.ServiceInfo{}
	// Inspect the AST to find the invocation of the wrapper function
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
					// If the function is the wrapper function, resolve the arguments for serviceName, Ip, and Port
					for _, arg := range n.Args {
						args = append(args, handleBasicLit(arg))
					}

					serviceName := resolveArgument(wrapper.ServiceName, args)
					ip := resolveArgument(wrapper.IP, args)
					port := resolveArgument(wrapper.Port, args)

					serviceInfos = append(serviceInfos, t.ServiceInfo{Application: service, IP: ip, Port: port})
					serviceNames = append(serviceNames, serviceName)
				}
			}
		}
		return true
	})

	return serviceNames, serviceInfos
}
