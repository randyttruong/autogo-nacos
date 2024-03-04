package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func ParseFile(filePath string) (*ast.File, error) {
	// ParseFile reads a Go source file and parses it into an abstract syntax tree (AST).
	//
	// filePath: The path to the Go source file.
	//
	// Returns:
	// A pointer to an ast.File struct representing the parsed Go source file.
	// An error if there was a problem reading the file or parsing the source code.

	// Convert file to an AST
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)

	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %v", filePath, err)
	}

	return fileAst, nil
}
