package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

func ParseYaml(filePath string) (*Yaml2Go, string, error) {
	conf := new(Yaml2Go)
	yamlFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
		return nil, "", err // Return the error
	}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
		return nil, "", err // Return the error
	}

	if conf.ApiVersion != "" && conf.Kind != ""  {
		// Return the configuration and the app label if present
		return conf, conf.Metadata.Name, nil
	}

	// If required fields are missing, return an error
	return nil, "", fmt.Errorf("missing required fields")
}