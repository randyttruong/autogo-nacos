package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	// "fmt"
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

func NewCustomVisitor(functionName string, wrapperFuncs []string) *customVisitor {
	return &customVisitor{functionName: functionName, wrapperFuncs: wrapperFuncs}
}

var wrapper = ""

func (v *customVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.FuncDecl:
		// Function declaration
		wrapper = n.Name.Name
		log.Printf("Function declaration: %s\n", n.Name.Name)

		// case *ast.AssignStmt:
		// 		// Variable assignment
		// 		for i, lhs := range n.Lhs {
	//     if len(n.Rhs) > i {
	//         switch rhs := n.Rhs[i].(type) {
	//         case *ast.BasicLit:
	//             // LHS is usually an Ident or a more complex expression.
	//             if ident, ok := lhs.(*ast.Ident); ok {
	//                 fmt.Printf("Assignment - Name: %s, Value: %s\n", ident.Name, rhs.Value)
	//             }
	//         }
	//     }
	// 		}
	// case *ast.ValueSpec:
	// 		// Variable/constant declaration
	// 		for i, name := range n.Names {
	//     // Assuming values are BasicLits for simplicity. Adjust as needed.
	//     if len(n.Values) > i {
	//         switch val := n.Values[i].(type) {
	//         case *ast.BasicLit:
	//             fmt.Printf("Var Declaration - Name: %s, Value: %s\n", name.Name, val.Value)
	//         }
	//     }
	// 		}
	case *ast.CallExpr:
		// Check if this callExpr matches the function of interest
		// For simplicity, let's assume we're looking for a function named "TargetFunction"
		if selExpr, ok := n.Fun.(*ast.SelectorExpr); ok {
			log.Printf("%s ", selExpr.Sel.Name)
			if selExpr.Sel.Name == v.functionName {
				log.Printf("%s called\n", v.functionName)
				log.Printf("%s is the wrapper\n", wrapper)
				// 				}
			}

			// Check the package and function name
			if selExpr.Sel.Name == "RegisterInstance" {
				// Found a call to TargetFunction
				// Analyze arguments here
				for _, arg := range n.Args {
					switch arg := arg.(type) {
					case *ast.BasicLit:
						log.Printf(arg.Value) // This is a literal value; you can access it via arg.Value
					case *ast.Ident:
						log.Printf(arg.Name)
                    case *ast.SelectorExpr:
                        log.Printf(arg.Sel.Name)
                    case *ast.CompositeLit:
                        if sel, ok := arg.Type.(*ast.SelectorExpr); ok {
                            log.Printf(sel.Sel.Name)
                            if sel.Sel.Name == "RegisterInstanceParam" {
                                for _, elt := range arg.Elts {
                                    if kv, ok := elt.(*ast.KeyValueExpr); ok {
                                        if key, ok := kv.Key.(*ast.Ident); ok {
                                            switch key.Name {
                                            case "Ip", "Port", "ServiceName":
                                                switch v := kv.Value.(type) {
                                                case *ast.Ident:
                                                    log.Printf("%s is a variable with value: %s", key.Name, v.Name)
                                                case *ast.BasicLit:
                                                    log.Printf("%s is a literal with value: %s", key.Name, v.Value)
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                                

                        }
                    }
				// }
			}
		}
	}}
	return v
}

// case *ast.CallExpr:
// 		// Function call
// 		if fun, ok := n.Fun.(*ast.SelectorExpr); ok {
// 				// Compare the function call to the specified function name
// 				if fun.Sel.Name == v.functionName {
// 						log.Printf("%s called\n", v.functionName)
// 				}
// 		}
// }

// if node == nil {
// 	return nil
// }
//   // Remove the unused variable declaration
//   var wrapper = ""

//   switch n := node.(type) {
//   // case *ast.ValueSpec:
// // 	for _, name := range n.Names {
// // 		v.varDeclarations[name.Name] = n
// // // 	}
// case *ast.FuncDecl:
//       wrapper = n.Name.Name
// 	if n.Type.Params != nil {
// 		for _, field := range n.Type.Params.List {
// 			for _, name := range field.Names {
// 				v.paramDeclarations[name.Name] = field
// 			}
// 		}
// 	}
// case *ast.CallExpr:
// 	// Check if this callExpr matches the function of interest
// 	// For simplicity, let's assume we're looking for a function named "TargetFunction"
// 	if selExpr, ok := n.Fun.(*ast.SelectorExpr); ok {
// 		if _, ok := selExpr.X.(*ast.Ident); ok {
// 			// Check the package and function name
// 			if selExpr.Sel.Name == "RegisterInstance" {
// 				// Found a call to TargetFunction
// 				// Analyze arguments here
//                   for _, arg := range n.Args {
//                       switch arg := arg.(type) {
//                       case *ast.BasicLit:
//                           print(arg.Value)// This is a literal value; you can access it via arg.Value
//                       case *ast.Ident:
//                           // This is an identifier; you need to trace it
//                           // Placeholder for tracing logic
//                       }
//                   }
// 			}
// 		}
// 	}
// }
// return v
// };
