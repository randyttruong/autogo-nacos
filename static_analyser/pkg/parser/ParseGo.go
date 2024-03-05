package parser

import (
	"fmt"
	"go/ast"
	t "static_analyser/pkg/types"
	"static_analyser/pkg/util"
	"strings"
)

// finds the wrappers for register instance
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

// finds the invocation of the wrappers for register instance and resolves the arguments for serviceName, Ip, and Port
func FindRegisterInstanceWrapperInvocations(node ast.Node, wrapper t.RegisterInstanceWrapper, service string) ([]string, []t.ServiceInfo) {
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
func FindServiceDiscoveryWrappers(node ast.Node) []t.ServiceDiscoveryWrapper {
	// Different nacos SDK functions for service discovery
	select_sdk := []string{"GetService", "SelectAllInstances", "SelectOneHealthyInstance", "SelectInstances", "Subscribe"}
	select_params := []string{"GetServiceParam", "SelectAllInstancesParam", "SelectOneHealthyInstanceParam", "SelectInstancesParam", "SubscribeParam"}

	var paramNames = []string{}
	var wrapper string
	var instances []t.ServiceDiscoveryWrapper
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
				if util.Contains(select_sdk, selExpr.Sel.Name) {
					// log.Printf("%s ", selExpr.Sel.Name)
					for _, arg := range n.Args {

						switch arg := arg.(type) {
						case *ast.CompositeLit:
							if sel, ok := arg.Type.(*ast.SelectorExpr); ok {

								if util.Contains(select_params, sel.Sel.Name) {

									instance := t.ServiceDiscoveryWrapper{}
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
																instance.ServiceName = util.FindConstValue(node, strings.TrimSpace(v.Name), wrapper)
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

func FindSelectInstanceWrappersInvocations(node ast.Node, wrapper t.ServiceDiscoveryWrapper, service string) []string {
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
