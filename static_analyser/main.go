package main

import (
	"fmt"
	"go/ast"
	// "io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var nacos_functions = []string{"RegisterInstance", "GetService", "SelectAllInstances", "SelectOneHealthyInstance", "SelectInstances", "Subscribe"}

func main() {
	var validYamlFiles []string
	parsedYamls := make(map[string]*Yaml2Go) // Initialize the map

	var service_to_manifest = make(map[string]TCPManifest)

	var occurrences = make(map[string][]string)

	var service_directory = make(map[string]string)

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
				service_directory[serviceName] = filepath.Dir(path)
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

	// Find all the go files with sdk function calls
	occurrences = FindGoFilesWithFunctions(root, nacos_functions)

	// Print the occurrences
	for funcName, files := range occurrences {
		fmt.Printf("Function %s is called in: %v\n", funcName, files)
	}

	log.Printf("\n\n\n\n\n\n\n")

	for func_name, files := range occurrences {
		var wrapper_funcs = []string{}

		for i, file := range files {
			fmt.Println(i, file)
			f := parseFile(file)

			if err != nil {
				log.Fatal(err)
			}
			visitor := NewCustomVisitor(func_name, wrapper_funcs)

			// Walk the AST with the visitor
			ast.Walk(visitor, f)
		}

	}

}
