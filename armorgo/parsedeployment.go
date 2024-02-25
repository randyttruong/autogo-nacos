package extractrequest
import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)


// Ports
type Ports struct {
	ContainerPort int `yaml:"containerPort"`
}

// TemplateMetadataLabels
type TemplateMetadataLabels struct {
	App string `yaml:"app"`
}

// TemplateSpec
type TemplateSpec struct {
	Containers []Containers `yaml:"containers"`
}

// Requests
type Requests struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// Limits
type Limits struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// LivenessProbe
type LivenessProbe struct {
	Exec LivenessProbeExec `yaml:"exec"`
}

// Labels
type Labels struct {
	App     string `yaml:"app"`
	Version string `yaml:"version"`
}

// Spec
type Spec struct {
	Template Template `yaml:"template"`
}

// TemplateMetadata
type TemplateMetadata struct {
	Labels TemplateMetadataLabels `yaml:"labels"`
}

// Containers
type Containers struct {
	Env            []Env          `yaml:"env"`
	Resources      Resources      `yaml:"resources"`
	Name           string         `yaml:"name"`
	Image          string         `yaml:"image"`
	Ports          []Ports        `yaml:"ports"`
	ReadinessProbe ReadinessProbe `yaml:"readinessProbe"`
	LivenessProbe  LivenessProbe  `yaml:"livenessProbe"`
}

// Yaml2Go
type Yaml2Go struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Env
type Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Exec
type Exec struct {
	Command []string `yaml:"command"`
}

// LivenessProbeExec
type LivenessProbeExec struct {
	Command []string `yaml:"command"`
}

// Metadata
type Metadata struct {
	Name   string `yaml:"name"`
	Labels Labels `yaml:"labels"`
}

// Template
type Template struct {
	Metadata TemplateMetadata `yaml:"metadata"`
	Spec     TemplateSpec     `yaml:"spec"`
}

// Resources
type Resources struct {
	Requests Requests `yaml:"requests"`
	Limits   Limits   `yaml:"limits"`
}

// ReadinessProbe
type ReadinessProbe struct {
	Exec Exec `yaml:"exec"`
}


func ParseYaml(filePath string) *Yaml2Go {
	// resultMap := make(map[string]interface{})
	conf := new(Yaml2Go)
	yamlFile, err := ioutil.ReadFile(filePath)


	//log.Println("yamlFile:", yamlFile)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	// err = yaml.Unmarshal(yamlFile, &resultMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	//log.Println("conf", conf.Spec.Template.Spec.Containers[0].Env)
	// log.Println("conf", resultMap)
	return conf
}



