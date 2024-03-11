package parser

import (
	"go/ast"
	t "static_analyser/pkg/types"
	"static_analyser/pkg/util"
	"strings"
)

func FindServiceDiscoveryWrappers(node ast.Node) []t.ServiceDiscoveryWrapper {
	// FindServiceDiscoveryWrappers traverses the AST (Abstract Syntax Tree) to find all instances of service discovery calls.
	//
	// node: The root node of the AST.
	//
	// Returns:
	// A slice of ServiceDiscoveryWrapper structs. Each struct represents a service discovery call found in the AST.
	// The ServiceDiscoveryWrapper struct contains the name of the wrapper function and the parameters passed to the service discovery call.

	select_sdk := []string{"GetService", "SelectAllInstances", "SelectOneHealthyInstance", "SelectInstances", "Subscribe"}
	select_params := []string{"GetServiceParam", "SelectAllInstancesParam", "SelectOneHealthyInstanceParam", "SelectInstancesParam", "SubscribeParam"}

	var paramNames = []string{}
	var wrapper string
	var instances []t.ServiceDiscoveryWrapper

	handleFuncDecl := func(n *ast.FuncDecl) {
		// handleFuncDecl is a closure that handles function declarations.
		//
		// n: The function declaration node.
		//
		// This closure sets the wrapper variable to the name of the function and adds the names of the parameters to the paramNames slice.

		wrapper = n.Name.Name
		for _, param := range n.Type.Params.List {
			for _, name := range param.Names {
				paramNames = append(paramNames, name.Name)
			}
		}
	}

	handleCallExpr := func(n *ast.CallExpr) {
		// handleCallExpr is a closure that handles call expressions.
		//
		// n: The call expression node.
		//
		// This closure checks if the call expression matches the criteria defined by select_sdk and select_params.
		// If it does, a new ServiceDiscoveryWrapper instance is created and added to the instances slice.

		selExpr, ok := n.Fun.(*ast.SelectorExpr)
		if !ok || !util.Contains(select_sdk, selExpr.Sel.Name) {
			return
		}

		for _, arg := range n.Args {
			arg, ok := arg.(*ast.CompositeLit)
			if !ok {
				continue
			}

			sel, ok := arg.Type.(*ast.SelectorExpr)
			if !ok || !util.Contains(select_params, sel.Sel.Name) {
				continue
			}

			instance := t.ServiceDiscoveryWrapper{Wrapper: wrapper}
			for _, elt := range arg.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}

				key, ok := kv.Key.(*ast.Ident)
				if !ok || key.Name != "ServiceName" {
					continue
				}

				switch v := kv.Value.(type) {
				case *ast.Ident:
					for i, paramName := range paramNames {
						if paramName == strings.TrimSpace(v.Name) {
							instance.ServiceName = t.WrapperParams{Position: i}
						}
						if instance.ServiceName == nil {
							instance.ServiceName = util.FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
						}
					}
				case *ast.BasicLit:
					instance.ServiceName = strings.ReplaceAll(strings.TrimSpace(v.Value), "\"", "")
				}
				instances = append(instances, instance)
			}
		}
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		case *ast.FuncDecl:
			handleFuncDecl(n)
		case *ast.CallExpr:
			handleCallExpr(n)
		}
		return true
	})
	return instances
}
