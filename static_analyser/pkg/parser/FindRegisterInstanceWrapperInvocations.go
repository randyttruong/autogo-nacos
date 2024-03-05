package parser

import (
	"fmt"
	"go/ast"
	t "static_analyser/pkg/types"
	"strings"
)

// finds the invocation of the wrappers for register instance and resolves the arguments for serviceName, Ip, and Port
func FindRegisterInstanceWrapperInvocations(node ast.Node, wrapper t.RegisterInstanceWrapper, service string) ([]string, []t.ServiceInfo) {

	handleBasicLit := func(arg ast.Expr) string {
		if lit, ok := arg.(*ast.BasicLit); ok {
			return lit.Value
		}
		return "nil"
	}

	resolveArgument := func(arg interface{}, args []string) string {
		switch t := arg.(type) {
		case string:
			return t
		case t.WrapperParams:
			return strings.ReplaceAll(args[t.Position], "\"", "")
		}
		return ""
	}

	wrapperName := wrapper.Wrapper
	serviceNames := []string{}
	serviceInfos := []t.ServiceInfo{}
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
						args = append(args, handleBasicLit(arg))
					}
					fmt.Printf("%v", args)

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

// // finds the invocation of the wrappers for register instance and resolves the arguments for serviceName, Ip, and Port
// func FindRegisterInstanceWrapperInvocations(node ast.Node, wrapper t.RegisterInstanceWrapper, service string) ([]string, []t.ServiceInfo) {
// 	wrapperName := wrapper.Wrapper
// 	serviceNames := []string{}
// 	serviceInfos := []t.ServiceInfo{}
// 	ast.Inspect(node, func(n ast.Node) bool {
// 		if n == nil {
// 			return false
// 		}

// 		switch n := n.(type) {

// 		case *ast.CallExpr:
// 			// log.Printf("Call expression: %s\n", n.Fun)
// 			if fun, ok := n.Fun.(*ast.Ident); ok {
// 				var args []string
// 				if fun.Name == wrapperName {

// 					for _, arg := range n.Args {

// 						if lit, ok := arg.(*ast.BasicLit); ok {
// 							args = append(args, lit.Value)
// 						} else {
// 							args = append(args, "nil")
// 						}
// 					}
// 					fmt.Printf("%v", args)

// 					serviceName := ""
// 					ip := ""
// 					port := ""

// 					switch t := wrapper.ServiceName.(type) {
// 					case string:
// 						serviceName = t
// 					case t.WrapperParams:
// 						serviceName = strings.ReplaceAll(args[t.Position], "\"", "")
// 					}
// 					switch t := wrapper.IP.(type) {
// 					case string:
// 						ip = t
// 					case t.WrapperParams:
// 						ip = strings.ReplaceAll(args[t.Position], "\"", "")
// 					}
// 					switch t := wrapper.Port.(type) {
// 					case string:
// 						port = t
// 					case t.WrapperParams:
// 						port = strings.ReplaceAll(args[t.Position], "\"", "")
// 					}
// 					serviceInfos = append(serviceInfos, t.ServiceInfo{Application: service, IP: ip, Port: port})
// 					serviceNames = append(serviceNames, serviceName)
// 				}

// 			}
// 			// var instances [].RegisterInfo

// 		}
// 		// return true
// 		return true
// 	})

// 	return serviceNames, serviceInfos
// }
