package parser

import (
	"go/ast"
	t "static_analyser/pkg/types"
	"static_analyser/pkg/util"
	"strings"
)

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
			// Function declaration
			wrapper = n.Name.Name
			paramNames = []string{}

			for _, param := range n.Type.Params.List {
				for _, name := range param.Names {
					paramNames = append(paramNames, name.Name)
				}
			}

		case *ast.CallExpr:
			if selExpr, ok := n.Fun.(*ast.SelectorExpr); ok {
				if selExpr.Sel.Name == "RegisterInstance" {
					for _, arg := range n.Args {
						switch arg := arg.(type) {
						case *ast.CompositeLit:
							if sel, ok := arg.Type.(*ast.SelectorExpr); ok {
								if sel.Sel.Name == "RegisterInstanceParam" {
									instance := t.RegisterInstanceWrapper{}
									instance.Wrapper = wrapper
									for _, elt := range arg.Elts {
										if kv, ok := elt.(*ast.KeyValueExpr); ok {
											if key, ok := kv.Key.(*ast.Ident); ok {

												switch key.Name {
												case "Ip", "Port", "ServiceName":
													switch v := kv.Value.(type) {
													case *ast.Ident:
														switch key.Name {
														case "Ip":
															for i, paramName := range paramNames {
																if paramName == strings.TrimSpace(v.Name) {
																	instance.IP = t.WrapperParams{Position: i}
																}
															}

															if instance.IP == nil {
																instance.IP = util.FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
															}

														case "Port":
															for i, paramName := range paramNames {

																if paramName == strings.TrimSpace(v.Name) {
																	instance.Port = t.WrapperParams{Position: i}
																}
																if instance.Port == nil {
																	instance.Port = util.FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
																}
															}

														case "ServiceName":
															for i, paramName := range paramNames {
																if paramName == strings.TrimSpace(v.Name) {
																	instance.ServiceName = t.WrapperParams{Position: i}
																}
																if instance.ServiceName == nil {
																	instance.ServiceName = util.FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
																}
															}
														}

													case *ast.BasicLit:
														switch key.Name {
														case "Ip":
															instance.IP = strings.TrimSpace(v.Value)
														case "Port":
															instance.Port = strings.TrimSpace(v.Value)
														case "ServiceName":
															instance.ServiceName = strings.TrimSpace(v.Value)
														}

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
				}
			}
		}
		return true
	})
	return instances
}
