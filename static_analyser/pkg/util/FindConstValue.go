package util

import (
	"go/ast"
)

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
