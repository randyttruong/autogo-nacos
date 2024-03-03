package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func parseFile(filePath string) *ast.File {
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Failed to parse file %s: %v", filePath, err)
	}
	return fileAst
}

type customVisitor struct {
	functionName string
	wrapperFuncs []string
}

// End implements ast.Node.
func (*customVisitor) End() token.Pos {
	panic("unimplemented")
}

// Pos implements ast.Node.
func (*customVisitor) Pos() token.Pos {
	panic("unimplemented")
}

func NewCustomVisitor(functionName string, wrapperFuncs []string) *customVisitor {
	return &customVisitor{functionName: functionName, wrapperFuncs: wrapperFuncs}
}

// finds the wrappers for register instance
func RegisterWrappers(node ast.Node) []RegisterInfo {
	var instances []RegisterInfo
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
									instance := RegisterInfo{}
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
																	instance.IP = WrapperParams{i}
																}
															}

														case "Port":
															for i, paramName := range paramNames {

																if paramName == strings.TrimSpace(v.Name) {
																	instance.Port = WrapperParams{i}
																}
															}

														case "ServiceName":
															for i, paramName := range paramNames {
																if paramName == strings.TrimSpace(v.Name) {
																	instance.ServiceName = WrapperParams{i}
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
func RegisterCalls(node ast.Node, wrapper RegisterInfo, service string) ([]string, []ServiceInfo) {
	wrapperName := wrapper.Wrapper
	serviceNames := []string{}
	serviceInfos := []ServiceInfo{}
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
					case WrapperParams:
						serviceName = strings.ReplaceAll(args[t.position], "\"", "")
					}
					switch t := wrapper.IP.(type) {
					case string:
						ip = t
					case WrapperParams:
						ip = strings.ReplaceAll(args[t.position], "\"", "")
					}
					switch t := wrapper.Port.(type) {
					case string:
						port = t
					case WrapperParams:
						port = strings.ReplaceAll(args[t.position], "\"", "")
					}
					serviceInfos = append(serviceInfos, ServiceInfo{service, ip, port})
					serviceNames = append(serviceNames, serviceName)
				}

			}
			// var instances []RegisterInfo

		}
		// return true
		return true
	})

	return serviceNames, serviceInfos
}

// finds the wrappers for service discovery
func DiscoveryWrappers(node ast.Node) []SelectInfo {
	// Different nacos SDK functions for service discovery
	select_sdk := []string{"GetService", "SelectAllInstances", "SelectOneHealthyInstance", "SelectInstances", "Subscribe"}
	select_params := []string{"GetServiceParam", "SelectAllInstancesParam", "SelectOneHealthyInstanceParam", "SelectInstancesParam", "SubscribeParam"}

	var paramNames = []string{}
	var wrapper string
	var instances []SelectInfo
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		case *ast.FuncDecl:
			// Function declaration
			wrapper = n.Name.Name
			log.Printf("Function declaration: %s\n", n.Name.Name)

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

									instance := SelectInfo{}
									instance.Wrapper = wrapper
									for _, elt := range arg.Elts {
										if kv, ok := elt.(*ast.KeyValueExpr); ok {
											if key, ok := kv.Key.(*ast.Ident); ok {
												if key.Name == "ServiceName" {
													switch v := kv.Value.(type) {
													case *ast.Ident:
														for i, paramName := range paramNames {
															if paramName == strings.TrimSpace(v.Name) {
																instance.ServiceName = WrapperParams{i}
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

func DiscoveryCalls(node ast.Node, wrapper SelectInfo, service string) []string {
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
					case WrapperParams:
						serviceName = strings.ReplaceAll(args[t.position], "\"", "")
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
