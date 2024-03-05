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
	fmt.Println("Valid .yaml files with required fields:")
	for _, file := range validYamlFiles {
		fmt.Println(file)
	}
	fmt.Println("")
}

func createTCPManifests(parsedYamls map[string]*t.Yaml2Go) map[string]t.TCPManifest {
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
	for application := range applicationFolders {
		temp := application2manifest[application]
		temp.Requests = callMap[application]
		fmt.Println("Manifest: ", temp)
		application2manifest[application] = temp

		f_util.WriteTCPManifestToJSON(application2manifest[application], application, outputPrefix)

	}
}

func main() {
	validYamlFiles, parsedYamls, applicationFolders, err := parseYamlFiles(root)
	if err != nil {
		fmt.Printf("Error walking the file tree: %v\n", err)
		return
	}

	printValidYamlFiles(validYamlFiles)

	// Create TCPManifest for each service
	application2manifest := createTCPManifests(parsedYamls)

	serviceDirectory, err := processServiceRegistrationCalls(applicationFolders, nacosFunctions)
	if err != nil {
		fmt.Printf("Error processing application folders: %v\n", err)
		return
	}

	callMap, err := processServiceDiscoveryCalls(applicationFolders, nacosFunctions, serviceDirectory)
	if err != nil {
		fmt.Printf("Error processing application files: %v\n", err)
		return
	}

	updateAndWriteManifests(applicationFolders, application2manifest, callMap, outputPrefix)

}
