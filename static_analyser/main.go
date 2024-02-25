package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)
var nacos_functions = []string{"RegisterInstance","GetService", "SelectAllInstances", "SelectOneHealthyInstance", "SelectInstances", "Subscribe"}

func main() {
	var validYamlFiles []string
	parsedYamls := make(map[string]*Yaml2Go) // Initialize the map

	var service_to_manifest = make(map[string]TCPManifest) 

	var occurrences = make(map[string][]string)

	root := "../example/" // Set the root directory you want to search

	// Walk thru all files in root dir and find yaml files with ApiVersion and Kind fields
	// Parse YAML files and store in map with service name
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			conf, serviceName, err := ParseYaml(path) // Correctly handle returned values
			if err != nil {
				fmt.Printf("Error parsing YAML file %s: %v\n", path, err)
				return nil // Continue processing other files even if this one fails
			}

			// Assuming you want to track YAML files that successfully parsed and contained the app label
			if serviceName != "" {
				parsedYamls[path] = conf
				validYamlFiles = append(validYamlFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the file tree: %v\n", err)
		return
	}

	fmt.Println("Valid .yaml files with required fields:")
	for _, file := range validYamlFiles {
		fmt.Println(file)
	}

	// Create TCPManifest for each service
	for service, value := range parsedYamls {
		version := value.Metadata.Labels.Version
	
		service_to_manifest[service] = TCPManifest{Version: version, Service: service}
	}
	
	// Output the TCPManifest for each service in JSON format
	for _, manifest := range service_to_manifest {
		PrintJson(manifest)
	}

// 	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
// 			fset := token.NewFileSet()
// 			node, err := parser.ParseFile(fset, path, nil, 0)
// 			if err != nil {
// 					fmt.Println(err)
// 					return err
// 			}

// 			ast.Inspect(node, func(n ast.Node) bool {
// 					// Check for function calls
// 					callExpr, ok := n.(*ast.CallExpr)
// 					if !ok {
// 							return true
// 					}
// 					switch fun := callExpr.Fun.(type) {
// 					case *ast.Ident: // Simple function calls
// 							if nacos_functions[fun.Name] {
// 									occurrences[fun.Name] = append(occurrences[fun.Name], path)
// 							}
// 					case *ast.SelectorExpr: // Qualified (package) function calls
// 							if id, ok := fun.X.(*ast.Ident); ok && nacos_functions[fun.Sel.Name] {
// 									fullName := id.Name + "." + fun.Sel.Name
// 									occurrences[fullName] = append(occurrences[fullName], path)
// 							}
// 					}
// 					return true
// 			})
// 	}
// 	return nil
// })
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
				return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
				content, err := ioutil.ReadFile(path)
				if err != nil {
						return err
				}
				contentStr := string(content)
				for _, funcName := range nacos_functions {
						if strings.Contains(contentStr, funcName+"(") { // Simple check for function call
								occurrences[funcName] = append(occurrences[funcName], path)
						}
				}
		}
		return nil
	})

	if err != nil {
			fmt.Printf("Error walking through the directory: %v\n", err)
	}

	// Print the occurrences
	for funcName, files := range occurrences {
			fmt.Printf("Function %s is called in: %v\n", funcName, files)
	}

	
}





