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

var nacosFunctions = []string{"RegisterInstance", "GetService", "SelectAllInstances", "SelectOneHealthyInstance", "SelectInstances", "Subscribe"}

func parseYamlFiles(root string) ([]string, map[string]*t.Yaml2Go, map[string]string, error) {
	// parseYamlFiles parses YAML files from a given root directory.
	//
	// root: The root directory for the search.
	//
	// Returns:
	// A list of paths to the valid YAML files.
	// A map where the keys are the names of the services and the values are pointers to the corresponding parsed YAML files.
	// A map where the keys are the names of the services and the values are the corresponding folder paths.
	// An error if there was a problem walking the file tree or parsing a YAML file.

	var validYamlFiles []string
	parsedYamls := make(map[string]*t.Yaml2Go) // Initialize the map
	var applicationFolders = make(map[string]string)

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
		return nil, nil, nil, err
	}

	return validYamlFiles, parsedYamls, applicationFolders, nil
}

func printValidYamlFiles(validYamlFiles []string) {
	// printValidYamlFiles prints the valid YAML files with required fields.
	//
	// validYamlFiles: A list of paths to the valid YAML files.
	//
	// Returns:
	// This function doesn't return a value. It prints the paths to the valid YAML files to the standard output.

	fmt.Println("Valid .yaml files with required fields:")
	for _, file := range validYamlFiles {
		fmt.Println(file)
	}
	fmt.Println("")
}

func createTCPManifests(parsedYamls map[string]*t.Yaml2Go) map[string]t.TCPManifest {
	// createTCPManifests creates TCPManifests from the parsed YAML files.
	//
	// parsedYamls: A map where the keys are the names of the applications and the values are pointers to the corresponding parsed YAML files.
	//
	// Returns:
	// A map where the keys are the names of the applications and the values are the corresponding TCPManifests.

	application2manifest := make(map[string]t.TCPManifest)

	for application, value := range parsedYamls {
		fmt.Printf("Service: %s, Version: %s \n", application, value.Metadata.Labels.Version)
		version := value.Metadata.Labels.Version
		application2manifest[application] = t.TCPManifest{Version: version, Service: application}
	}
	fmt.Println("")

	return application2manifest
}

func processServiceRegistrationCalls(applicationFolders map[string]string, nacosFunctions []string) (map[string]t.ServiceInfo, error) {
	// processServiceRegistrationCalls processes the service registration calls from the application folders.
	//
	// applicationFolders: A map where the keys are the names of the applications and the values are the corresponding folder paths.
	// nacosFunctions: A list of function names to search for in the .go files.
	//
	// Returns:
	// A map where the keys are the names of the services and the values are the corresponding ServiceInfo.
	// An error if there was a problem finding go files, parsing a file, or finding service registration wrappers.

	registerWrapperMap := make(map[string]t.RegisterInstanceWrapper)
	serviceDirectory := make(map[string]t.ServiceInfo)
	var nacosFiles map[string][]string

	for application, dir := range applicationFolders {
		// Find all the go files
		goFiles, err := file_finder.FindGoFiles(dir)
		if err != nil {
			return nil, fmt.Errorf("error finding go files in %s: %v", dir, err)
		}
		// Find all the go files with sdk function calls
		nacosFiles, err = file_finder.FindGoFilesWithFunctions(dir, nacosFunctions)
		if err != nil {
			return nil, fmt.Errorf("error finding go files with nacos functions in %s: %v", dir, err)
		}
		for _, files := range nacosFiles {
			for _, file := range files {
				f, err := parser.ParseFile(file)
				if err != nil {
					return nil, fmt.Errorf("error parsing file %s: %v", file, err)
				}
				instances := parser.FindRegisterInstanceWrappers(f)
				for _, instance := range instances {
					registerWrapperMap[application] = instance
				}
			}
			for _, file := range goFiles {
				f, err := parser.ParseFile(file)
				if err != nil {
					return nil, fmt.Errorf("error parsing file %s: %v", file, err)
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
	return serviceDirectory, nil
}

func processServiceDiscoveryCalls(applicationFolders map[string]string, nacosFunctions []string, serviceDirectory map[string]t.ServiceInfo) (map[string][]t.TCPRequest, error) {
	// processServiceDiscoveryCalls processes the service discovery calls from the application folders.
	//
	// applicationFolders: A map where the keys are the names of the applications and the values are the corresponding folder paths.
	// nacosFunctions: A list of function names to search for in the .go files.
	// serviceDirectory: A map where the keys are the names of the services and the values are the corresponding ServiceInfo.
	//
	// Returns:
	// A map where the keys are the names of the applications and the values are slices of TCPRequests.
	// An error if there was a problem finding go files, parsing a file, or finding service discovery wrappers.

	serviceDiscoveryWrapperMap := make(map[string]t.ServiceDiscoveryWrapper)
	callMap := make(map[string][]t.TCPRequest)
	var nacosFiles map[string][]string

	for application, dir := range applicationFolders {
		goFiles, err := file_finder.FindGoFiles(dir)
		if err != nil {
			return nil, fmt.Errorf("error finding go files in %s: %v", dir, err)
		}

		nacosFiles, err = file_finder.FindGoFilesWithFunctions(dir, nacosFunctions)
		if err != nil {
			return nil, fmt.Errorf("error finding go files with nacos functions in %s: %v", dir, err)
		}

		for _, files := range nacosFiles {
			for _, file := range files {
				f, err := parser.ParseFile(file)
				if err != nil {
					return nil, fmt.Errorf("error parsing file %s: %v", file, err)
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
				return nil, fmt.Errorf("error parsing file %s: %v", file, err)
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

	return callMap, nil
}

func updateAndWriteManifests(applicationFolders map[string]string, application2manifest map[string]t.TCPManifest, callMap map[string][]t.TCPRequest, outputPrefix string) {
	// updateAndWriteManifests updates the TCPManifests with the corresponding TCPRequests and writes them to JSON files.
	//
	// applicationFolders: A map where the keys are the names of the applications and the values are the corresponding folder paths.
	// application2manifest: A map where the keys are the names of the applications and the values are the corresponding TCPManifests.
	// callMap: A map where the keys are the names of the applications and the values are slices of TCPRequests.
	// outputPrefix: A string to be prepended to the output file names.
	//
	// Returns:
	// This function doesn't return a value. It updates the TCPManifests in application2manifest and writes them to JSON files. The file names are the application names with outputPrefix prepended.

	for application := range applicationFolders {
		temp := application2manifest[application]
		temp.Requests = callMap[application]
		fmt.Println("Manifest: ", temp)
		application2manifest[application] = temp

		f_util.WriteTCPManifestToJSON(application2manifest[application], application, outputPrefix)

	}
}

func main() {
	// main is the entry point of the program.
	// It performs the following steps:
	// 1. Parses YAML files from a given root directory.
	// 2. Prints the valid YAML files.
	// 3. Creates TCP manifests from the parsed YAMLs.
	// 4. Processes service registration calls from the application folders.
	// 5. Processes service discovery calls from the application folders.
	// 6. Updates and writes the manifests.

	// Parse YAML files from the root directory
	validYamlFiles, parsedYamls, applicationFolders, err := parseYamlFiles(root)
	if err != nil {
		fmt.Printf("Error walking the file tree: %v\n", err)
		return
	}

	// Print the valid YAML files
	printValidYamlFiles(validYamlFiles)

	// Create TCP manifests from the parsed YAMLs
	application2manifest := createTCPManifests(parsedYamls)

	// Process service registration calls from the application folders
	serviceDirectory, err := processServiceRegistrationCalls(applicationFolders, nacosFunctions)
	if err != nil {
		fmt.Printf("Error processing application folders: %v\n", err)
		return
	}

	// Process service discovery calls from the application folders
	callMap, err := processServiceDiscoveryCalls(applicationFolders, nacosFunctions, serviceDirectory)
	if err != nil {
		fmt.Printf("Error processing application files: %v\n", err)
		return
	}

	// Update and write the manifests
	updateAndWriteManifests(applicationFolders, application2manifest, callMap, outputPrefix)
}
