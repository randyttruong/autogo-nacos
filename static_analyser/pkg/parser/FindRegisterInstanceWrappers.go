package parser

import (
	"go/ast"
	t "static_analyser/pkg/types"
	"static_analyser/pkg/util"
	"strings"
)

func handleFuncDecl(n *ast.FuncDecl) (string, []string) {
	wrapper := n.Name.Name
	paramNames := []string{}
	for _, param := range n.Type.Params.List {
		for _, name := range param.Names {
			paramNames = append(paramNames, name.Name)
		}
	}
	return wrapper, paramNames
}

func handleIdent(v *ast.Ident, keyName string, paramNames []string, instance t.RegisterInstanceWrapper, node ast.Node, wrapper string) t.RegisterInstanceWrapper {
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

func handleBasicLit(v *ast.BasicLit, keyName string, instance t.RegisterInstanceWrapper) t.RegisterInstanceWrapper {
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

func handleCallExpr(n *ast.CallExpr, node ast.Node, wrapper string, paramNames []string, instances []t.RegisterInstanceWrapper) []t.RegisterInstanceWrapper {
	if selExpr, ok := n.Fun.(*ast.SelectorExpr); ok && selExpr.Sel.Name == "RegisterInstance" {
		for _, arg := range n.Args {
			switch arg := arg.(type) {
			case *ast.CompositeLit:
				if sel, ok := arg.Type.(*ast.SelectorExpr); ok && sel.Sel.Name == "RegisterInstanceParam" {
					instance := t.RegisterInstanceWrapper{}
					instance.Wrapper = wrapper
					for _, elt := range arg.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok {
								switch key.Name {
								case "Ip", "Port", "ServiceName":
									switch v := kv.Value.(type) {
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

func FindRegisterInstanceWrappers(node ast.Node) []t.RegisterInstanceWrapper {
	var instances []t.RegisterInstanceWrapper
	var paramNames = []string{}
	var wrapper string
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
