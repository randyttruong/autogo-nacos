package types

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

type TCPManifest struct {
	Service  string       `json:"service"`
	Version  string       `json:"version"`
	Requests []TCPRequest `json:"requests"`
}

type TCPRequest struct {
	Type string `json:"type"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Port string `json:"port"`
}

type endpoint struct {
	name string
	path string
}

type tcpEndpoint struct {
	//pkgName string
	name string
	//path    string
	port string
}

// Stores information about wrappers for registering
type RegisterInfo struct {
	Wrapper string
	// Interface is string if varaible is hard cades and it is WrapperParams if variable is a argument to the wrapper
	ServiceName interface{}
	IP          interface{}
	Port        interface{}
}

type WrapperParams struct {
	// position the argument is passed into the wrapper
	Position int
}

// Stores information about wrappers for selecting
type SelectInfo struct {
	Wrapper     string
	ServiceName interface{}
}

// stores information about the service
type ServiceInfo struct {
	Application string
	IP          string
	Port        string
}
