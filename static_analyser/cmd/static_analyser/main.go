package main

import (
	"fmt"
	"os"
	"path/filepath"
	f_util "static_analyser/pkg/fileUtils"
	"static_analyser/pkg/file_finder"
	"static_analyser/pkg/parser"
	t "static_analyser/pkg/types"
	"strings"
)

// set the output directory for manifests
var outputPrefix = "../output/game_microservices/"

// var outputPrefix = "output\\game_microservices\\"

// Set the root directory you want to search
var root = "../tests/example_2/"

// var root = "../input/"
// var root = "..\\input\\"

var nacos_functions = []string{"RegisterInstance", "GetService", "SelectAllInstances", "SelectOneHealthyInstance", "SelectInstances", "Subscribe"}

func updateAndWriteManifests(applicationFolders map[string]string, application2manifest map[string]t.TCPManifest, callMap map[string][]t.TCPRequest, outputPrefix string) {
	for application := range applicationFolders {
		temp := application2manifest[application]
		temp.Requests = callMap[application]
		fmt.Println("Manifest: ", temp)
		application2manifest[application] = temp

		f_util.WriteTCPManifestToJSON(application2manifest[application], application, outputPrefix)

	}
}

func main() {
	var validYamlFiles []string
	parsedYamls := make(map[string]*t.Yaml2Go) // Initialize the map

	// map to store the TCPManifest for each application
	var application2manifest = make(map[string]t.TCPManifest)

	// store the files with nacos functions for each application
	var nacosFiles = make(map[string][]string)

	// store the folders with containing the code for each application
	var applicationFolders = make(map[string]string)

	//store the information for each service registrated to nacos for each application
	// maps application to service information
	var serviceDirectory = make(map[string]t.ServiceInfo)

	// store the service discovery calls to nacos for each application
	var callMap = make(map[string][]t.TCPRequest)

	var registerWrapperMap = make(map[string]t.RegisterInstanceWrapper)
	var serviceDiscoveryWrapperMap = make(map[string]t.ServiceDiscoveryWrapper)

	// Walk thru all files in root dir and find yaml files with ApiVersion and Kind fields
	// Parse YAML files and store in map with service name
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			conf, serviceName, err := parser.ParseYaml(path) // Correctly handle returned values
			if err != nil {
				fmt.Printf("Error parsing YAML file %s: %v\n", path, err)
				return nil // Continue processing other files even if this one fails
			}

			// Assuming you want to track YAML files that successfully parsed and contained the app label
			if serviceName != "" {
				parsedYamls[serviceName] = conf
				applicationFolders[serviceName] = filepath.Dir(path)
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
	fmt.Println("")

	// Create TCPManifest for each service
	for application, value := range parsedYamls {
		fmt.Printf("Service: %s, Version: %s \n", application, value.Metadata.Labels.Version)
		version := value.Metadata.Labels.Version
		application2manifest[application] = t.TCPManifest{Version: version, Service: application}
	}
	fmt.Println("")

	// Only looks at the directories corresponding to each service.
	// Loop for registration  calls
	for application, dir := range applicationFolders {
		// Find all the go files
		goFiles, err := file_finder.FindGoFiles(dir)

		if err != nil {
			fmt.Println("lol")
			fmt.Printf("error finding go files in %s: %v\n", dir, err)
			return
		}
		// Find all the go files with sdk function calls
		nacosFiles, err = file_finder.FindGoFilesWithFunctions(dir, nacos_functions)

		if err != nil {
			fmt.Printf("error finding go files with nacos functions in %s: %v\n", dir, err)
			return
		}
		for _, files := range nacosFiles {

			for _, file := range files {
				f, err := parser.ParseFile(file)

				if err != nil {
					fmt.Printf("error parsing file %s: %v\n", file, err)
					return
				}

				instances := parser.FindRegisterInstanceWrappers(f)
				for _, instance := range instances {
					registerWrapperMap[application] = instance
				}
			}

			for _, file := range goFiles {
				f, err := parser.ParseFile(file)

				if err != nil {
					fmt.Printf("error parsing file %s: %v\n", file, err)
					return
				}
				for key, value := range registerWrapperMap {
					if key == application {
						names, infos := parser.FindRegisterInstanceWrapperInvocations(f, value, application)
						for i, name := range names {
							serviceDirectory[name] = infos[i]
						}
					}
				}
			}
		}
	}

	// Loop for service discovery calls
	for application, dir := range applicationFolders {
		goFiles, err := file_finder.FindGoFiles(dir)

		if err != nil {
			fmt.Printf("error finding go files in %s: %v\n", dir, err)
			return
		}
		// Find all the go files with sdk function calls
		nacosFiles, err = file_finder.FindGoFilesWithFunctions(dir, nacos_functions)

		if err != nil {
			fmt.Printf("error finding go files with nacos functions in %s: %v\n", dir, err)
			return
		}

		for _, files := range nacosFiles {

			for _, file := range files {
				f, err := parser.ParseFile(file)

				if err != nil {
					fmt.Printf("error parsing file %s: %v\n", file, err)
					return
				}

				instances := parser.FindServiceDiscoveryWrappers(f)
				for _, instance := range instances {
					serviceDiscoveryWrapperMap[application] = instance

				}
			}
		}

		for _, file := range goFiles {

			f, err := parser.ParseFile(file)

			if err != nil {
				fmt.Printf("error parsing file %s: %v\n", file, err)
				return
			}

			for key, value := range serviceDiscoveryWrapperMap {
				if key == application {
					names := parser.FindSelectInstanceWrappersInvocations(f, value, application)
					for _, name := range names {
						req := t.TCPRequest{Type: "tcp", URL: serviceDirectory[name].IP, Name: serviceDirectory[name].Application, Port: serviceDirectory[name].Port}
						callMap[application] = append(callMap[application], req)
					}

				}

			}
		}
	}

	updateAndWriteManifests(applicationFolders, application2manifest, callMap, outputPrefix)

}
