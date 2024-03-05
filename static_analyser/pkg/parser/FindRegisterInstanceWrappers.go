package parser

import (
	"go/ast"
	t "static_analyser/pkg/types"
	"static_analyser/pkg/util"
	"strings"
)

func FindRegisterInstanceWrappers(node ast.Node) []t.RegisterInstanceWrapper {
	// FindRegisterInstanceWrappers traverses the AST (Abstract Syntax Tree) to find all instances of RegisterInstance calls.
	//
	// node: The root node of the AST.
	//
	// Returns:
	// A slice of RegisterInstanceWrapper structs. Each struct represents a RegisterInstance call found in the AST.
	// The RegisterInstanceWrapper struct contains the name of the wrapper function and the parameters passed to the RegisterInstance call.

	handleFuncDecl := func(n *ast.FuncDecl) (string, []string) {
		wrapper := n.Name.Name
		paramNames := []string{}

		// Iterate over the parameters of the function
		for _, param := range n.Type.Params.List {
			for _, name := range param.Names {
				paramNames = append(paramNames, name.Name)
			}
		}
		return wrapper, paramNames
	}

	handleIdent := func(v *ast.Ident, keyName string, paramNames []string, instance t.RegisterInstanceWrapper, node ast.Node, wrapper string) t.RegisterInstanceWrapper {
		// handleIdent processes an *ast.Ident node and updates the corresponding field in the given RegisterInstanceWrapper struct.
		//
		// v: The *ast.Ident node to process.
		// keyName: The name of the field to update in the RegisterInstanceWrapper struct. It should be one of "Ip", "Port", or "ServiceName".
		// paramNames: A slice of parameter names from the wrapper function.
		// instance: The RegisterInstanceWrapper struct to update.
		// node: The root node of the AST.
		// wrapper: The name of the wrapper function.
		//
		// Returns:
		// The updated RegisterInstanceWrapper struct.

		// Check if the Ident is a parameter of the function
		for i, paramName := range paramNames {
			if paramName == strings.TrimSpace(v.Name) {
				switch keyName {
				case "Ip":
					instance.IP = t.WrapperParams{Position: i}
				case "Port":
					instance.Port = t.WrapperParams{Position: i}
				case "ServiceName":
					instance.ServiceName = t.WrapperParams{Position: i}
				}
			}
		}
		if instance.IP == nil || instance.Port == nil || instance.ServiceName == nil {
			switch keyName {
			case "Ip":
				instance.IP = util.FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
			case "Port":
				instance.Port = util.FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
			case "ServiceName":
				instance.ServiceName = util.FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
			}
		}
		return instance
	}

	handleBasicLit := func(v *ast.BasicLit, keyName string, instance t.RegisterInstanceWrapper) t.RegisterInstanceWrapper {
		// handleBasicLit processes an *ast.BasicLit node and updates the corresponding field in the given RegisterInstanceWrapper struct.
		//
		// v: The *ast.BasicLit node to process.
		// keyName: The name of the field to update in the RegisterInstanceWrapper struct. It should be one of "Ip", "Port", or "ServiceName".
		// instance: The RegisterInstanceWrapper struct to update.
		//
		// Returns:
		// The updated RegisterInstanceWrapper struct.

		switch keyName {
		case "Ip":
			instance.IP = strings.TrimSpace(v.Value)
		case "Port":
			instance.Port = strings.TrimSpace(v.Value)
		case "ServiceName":
			instance.ServiceName = strings.TrimSpace(v.Value)
		}
		return instance
	}

	handleCallExpr := func(n *ast.CallExpr, node ast.Node, wrapper string, paramNames []string, instances []t.RegisterInstanceWrapper) []t.RegisterInstanceWrapper {
		// handleCallExpr processes an *ast.CallExpr node to find instances of RegisterInstance calls.
		//
		// n: The *ast.CallExpr node to process.
		// node: The root node of the AST.
		// wrapper: The name of the wrapper function.
		// paramNames: A slice of parameter names from the wrapper function.
		// instances: A slice of RegisterInstanceWrapper structs found so far.
		//
		// Returns:
		// A slice of RegisterInstanceWrapper structs. If the *ast.CallExpr node represents a RegisterInstance call, a new RegisterInstanceWrapper struct is created and added to the slice.
		//

		// Check if the function is RegisterInstance
		if selExpr, ok := n.Fun.(*ast.SelectorExpr); ok && selExpr.Sel.Name == "RegisterInstance" {
			for _, arg := range n.Args {
				switch arg := arg.(type) {
				// Check if the argument is a CompositeLit
				case *ast.CompositeLit:
					// Check if the CompositeLit is of type RegisterInstanceParam
					if sel, ok := arg.Type.(*ast.SelectorExpr); ok && sel.Sel.Name == "RegisterInstanceParam" {
						instance := t.RegisterInstanceWrapper{}
						instance.Wrapper = wrapper
						// Iterate over the elements of the CompositeLit
						for _, elt := range arg.Elts {
							// Check if the element is a KeyValueExpr
							if kv, ok := elt.(*ast.KeyValueExpr); ok {
								// Check if the key is an Identifier
								if key, ok := kv.Key.(*ast.Ident); ok {
									switch key.Name {
									// Check if the key is Ip, Port or ServiceName
									case "Ip", "Port", "ServiceName":
										switch v := kv.Value.(type) {
										// Check if the value is an Ident or BasicLit
										case *ast.Ident:
											instance = handleIdent(v, key.Name, paramNames, instance, node, wrapper)
										case *ast.BasicLit:
											instance = handleBasicLit(v, key.Name, instance)
										}
									}
								}
							}
						}
						instances = append(instances, instance)
					}
				}
			}
		}
		return instances
	}

	var instances []t.RegisterInstanceWrapper
	var paramNames = []string{}
	var wrapper string

	// Inspect the AST and look for RegisterInstance calls
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		case *ast.FuncDecl:
			wrapper, paramNames = handleFuncDecl(n)

		case *ast.CallExpr:
			instances = handleCallExpr(n, node, wrapper, paramNames, instances)
		}

		return true
	})

	return instances
}
