package parser

import (
	"fmt"
	"go/ast"
	t "static_analyser/pkg/types"
	"strings"
)

// finds the wrappers for register instance
func RegisterWrappers(node ast.Node) []t.RegisterInfo {
	var instances []t.RegisterInfo
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
									instance := t.RegisterInfo{}
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
																instance.IP = FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
															}

														case "Port":
															for i, paramName := range paramNames {

																if paramName == strings.TrimSpace(v.Name) {
																	instance.Port = t.WrapperParams{Position: i}
																}
																if instance.Port == nil {
																	instance.Port = FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
																}
															}

														case "ServiceName":
															for i, paramName := range paramNames {
																if paramName == strings.TrimSpace(v.Name) {
																	instance.ServiceName = t.WrapperParams{Position: i}
																}
																if instance.ServiceName == nil {
																	instance.ServiceName = FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
																}
															}
															// if instance.ServiceName == nil {
															// 	instance.ServiceName = v.Name
															// }
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

// finds the invocation of the wrappers for register instance and resolves the arguments for serviceName, Ip, and Port
func RegisterCalls(node ast.Node, wrapper t.RegisterInfo, service string) ([]string, []t.ServiceInfo) {
	wrapperName := wrapper.Wrapper
	serviceNames := []string{}
	serviceInfos := []t.ServiceInfo{}
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {

		case *ast.CallExpr:
			// log.Printf("Call expression: %s\n", n.Fun)
			if fun, ok := n.Fun.(*ast.Ident); ok {
				var args []string
				if fun.Name == wrapperName {

					for _, arg := range n.Args {

						if lit, ok := arg.(*ast.BasicLit); ok {
							args = append(args, lit.Value)
						} else {
							args = append(args, "nil")
						}
					}
					fmt.Printf("%v", args)

					serviceName := ""
					ip := ""
					port := ""

					switch t := wrapper.ServiceName.(type) {
					case string:
						serviceName = t
					case t.WrapperParams:
						serviceName = strings.ReplaceAll(args[t.Position], "\"", "")
					}
					switch t := wrapper.IP.(type) {
					case string:
						ip = t
					case t.WrapperParams:
						ip = strings.ReplaceAll(args[t.Position], "\"", "")
					}
					switch t := wrapper.Port.(type) {
					case string:
						port = t
					case t.WrapperParams:
						port = strings.ReplaceAll(args[t.Position], "\"", "")
					}
					serviceInfos = append(serviceInfos, t.ServiceInfo{Application: service, IP: ip, Port: port})
					serviceNames = append(serviceNames, serviceName)
				}

			}
			// var instances [].RegisterInfo

		}
		// return true
		return true
	})

	return serviceNames, serviceInfos
}

// finds the wrappers for service discovery
func DiscoveryWrappers(node ast.Node) []t.SelectInfo {
	// Different nacos SDK functions for service discovery
	select_sdk := []string{"GetService", "SelectAllInstances", "SelectOneHealthyInstance", "SelectInstances", "Subscribe"}
	select_params := []string{"GetServiceParam", "SelectAllInstancesParam", "SelectOneHealthyInstanceParam", "SelectInstancesParam", "SubscribeParam"}

	var paramNames = []string{}
	var wrapper string
	var instances []t.SelectInfo
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		case *ast.FuncDecl:
			// Function declaration
			wrapper = n.Name.Name
			for _, param := range n.Type.Params.List {
				for _, name := range param.Names {
					paramNames = append(paramNames, name.Name)
				}
			}
			// log.Printf("Parameter names: %v\n", paramNames)

		case *ast.CallExpr:
			// log.Printf("Call expression: %s\n", n.Fun)
			if selExpr, ok := n.Fun.(*ast.SelectorExpr); ok {
				// If the function is a list of nacos sdk functions
				if contains(select_sdk, selExpr.Sel.Name) {
					// log.Printf("%s ", selExpr.Sel.Name)
					for _, arg := range n.Args {

						switch arg := arg.(type) {
						case *ast.CompositeLit:
							if sel, ok := arg.Type.(*ast.SelectorExpr); ok {

								if contains(select_params, sel.Sel.Name) {

									instance := t.SelectInfo{}
									instance.Wrapper = wrapper
									for _, elt := range arg.Elts {
										if kv, ok := elt.(*ast.KeyValueExpr); ok {
											if key, ok := kv.Key.(*ast.Ident); ok {
												if key.Name == "ServiceName" {
													switch v := kv.Value.(type) {
													case *ast.Ident:
														for i, paramName := range paramNames {
															if paramName == strings.TrimSpace(v.Name) {
																instance.ServiceName = t.WrapperParams{Position: i}
															}
															if instance.ServiceName == nil {
																instance.ServiceName = FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
															}
														}
													case *ast.BasicLit:
														instance.ServiceName = strings.ReplaceAll(strings.TrimSpace(v.Value), "\"", "")
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

func DiscoveryCalls(node ast.Node, wrapper t.SelectInfo, service string) []string {
	wrapperName := wrapper.Wrapper
	serviceNames := []string{}

	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		case *ast.CallExpr:
			if fun, ok := n.Fun.(*ast.Ident); ok {
				var args []string
				if fun.Name == wrapperName {

					for _, arg := range n.Args {

						if lit, ok := arg.(*ast.BasicLit); ok {
							args = append(args, lit.Value)
						} else {
							args = append(args, "nil")
						}
					}
					fmt.Printf("%v", args)

					serviceName := ""

					switch t := wrapper.ServiceName.(type) {
					case string:
						serviceName = t
					case t.WrapperParams:
						serviceName = strings.ReplaceAll(args[t.Position], "\"", "")
					}

					serviceNames = append(serviceNames, serviceName)
				}

			}

		}

		return true
	})

	return serviceNames
}

// helper for checking if a string is in a slice
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func FindConstValue(root ast.Node, constName string, wrapper string) string {
	var constValue string
	curr_wrapper := ""
	ast.Inspect(root, func(n ast.Node) bool {

		switch node := n.(type) {
		case *ast.FuncDecl:
			curr_wrapper = node.Name.Name

		case *ast.AssignStmt:
			for _, lhs := range node.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok && ident.Name == constName && curr_wrapper == wrapper {
					if len(node.Rhs) > 0 {
						if rhs, ok := node.Rhs[0].(*ast.BasicLit); ok {
							constValue = rhs.Value
							return false
						}
						if callExpr, ok := node.Rhs[0].(*ast.CallExpr); ok {
							if len(callExpr.Args) > 0 {
								if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
									constValue = basicLit.Value
									return false
								}
							}
						}
					}
				}
			}

		}
		return true // continue the inspection otherwise
	})
	return constValue
}
