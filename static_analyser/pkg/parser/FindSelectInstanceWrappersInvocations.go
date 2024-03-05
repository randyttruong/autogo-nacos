package parser

import (
	"fmt"
	"go/ast"
	t "static_analyser/pkg/types"
	"strings"
)

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
