package util

import (
	"go/ast"
)

func FindConstValue(root ast.Node, constName string, wrapper string) string {
	// FindConstValue inspects the AST (Abstract Syntax Tree) to find the value of a specific constant.
	//
	// root: The root node of the AST.
	// constName: The name of the constant whose value is to be found.
	// wrapper: The name of the wrapper function.
	//
	// Returns:
	// The value of the constant if found. Returns an empty string if the constant is not found.

	var constValue string
	curr_wrapper := ""

	// Begin analyzing the AST
	ast.Inspect(root, func(n ast.Node) bool {

		// Check the type of the node
		switch node := n.(type) {

		// Check if the node is a function declaration
		case *ast.FuncDecl:
			curr_wrapper = node.Name.Name

		// Check if the node is an assignment statement
		case *ast.AssignStmt:
			// Iterate through the left-hand side of the assignment
			for _, lhs := range node.Lhs {
				// Check if the left-hand side is an identifier with the specified name
				if ident, ok := lhs.(*ast.Ident); !ok || ident.Name != constName || curr_wrapper != wrapper {
					continue
				}
				// Check if there is a  right-hand side of the assignment
				if len(node.Rhs) <= 0 {
					continue
				}
				// Check if the right-hand side of the assignment is a basic literal
				if rhs, ok := node.Rhs[0].(*ast.BasicLit); ok {
					constValue = rhs.Value
					return false
				}

				// Check if the right-hand side of the assignment is a call expression
				if callExpr, ok := node.Rhs[0].(*ast.CallExpr); ok && len(callExpr.Args) > 0 {
					// Check if the argument of the call expression is a basic literal
					if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
						constValue = basicLit.Value
						return false
					}
				}

			}
		}
		return true // continue the inspection otherwise
	})
	return constValue
}
