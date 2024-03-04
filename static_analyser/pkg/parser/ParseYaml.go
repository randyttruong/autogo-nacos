package parser

import (
	"fmt"
	"os"
	t "static_analyser/pkg/types"

	"gopkg.in/yaml.v2"
)

func ParseYaml(filePath string) (*t.Yaml2Go, string, error) {
	// ParseYaml reads a YAML file and unmarshals it into a Yaml2Go struct.
	//
	// filePath: The path to the YAML file.
	//
	// Returns:
	// A pointer to a Yaml2Go struct containing the unmarshaled data.
	// The name of the metadata as a string.
	// An error if there was a problem reading the file, unmarshaling the data, or if required fields are missing.

	conf := new(t.Yaml2Go)

	// Read the file
	yamlFile, err := os.ReadFile(filePath)
	if err != nil{
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the YAML file into the configuration struct
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal YAML: %w", err)
  }

	// Check if the required fields are present
	if conf.ApiVersion == "" && conf.Kind == ""  {
		return nil, "", fmt.Errorf("missing required fields")
	}

	return conf, conf.Metadata.Name, nil
}
