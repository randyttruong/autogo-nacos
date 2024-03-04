package types

// Containers represents a collection of containers in a configuration.
type Containers struct {
	Env            []Env          `yaml:"env"`            // Environment variables for the containers.
	Resources      Resources      `yaml:"resources"`      // Resource limits and requests for the containers.
	Name           string         `yaml:"name"`           // Name of the container.
	Image          string         `yaml:"image"`          // Docker image for the container.
	Ports          []Ports        `yaml:"ports"`          // Ports to expose for the container.
	ReadinessProbe ReadinessProbe `yaml:"readinessProbe"` // Configuration for the readiness probe.
	LivenessProbe  LivenessProbe  `yaml:"livenessProbe"`  // Configuration for the liveness probe.
}

// Env represents an environment variable.
type Env struct {
	Name  string `yaml:"name"`  // Name is the name of the environment variable.
	Value string `yaml:"value"` // Value is the value of the environment variable.
}

// Exec represents a command execution configuration.
type Exec struct {
	Command []string `yaml:"command"`
}

// Labels represents the labels associated with a resource.
type Labels struct {
	App     string `yaml:"app"`     // App represents the application name.
	Version string `yaml:"version"` // Version represents the version of the resource.
}

// Limits represents the resource limits for a particular task.
type Limits struct {
	Cpu    string `yaml:"cpu"`    // Cpu represents the CPU limit for the task.
	Memory string `yaml:"memory"` // Memory represents the memory limit for the task.
}

// LivenessProbe represents the liveness probe configuration for a service.
type LivenessProbe struct {
	Exec LivenessProbeExec `yaml:"exec"`
}

// LivenessProbeExec represents the configuration for executing a command as a liveness probe.
type LivenessProbeExec struct {
	Command []string `yaml:"command"`
}

// Metadata represents the metadata associated with a resource.
type Metadata struct {
	Name   string `yaml:"name"`   // Name is the name of the resource.
	Labels Labels `yaml:"labels"` // Labels are the labels associated with the resource.
}

// Ports represents the ports configuration for a container.
type Ports struct {
	ContainerPort int `yaml:"containerPort"`
}

// ReadinessProbe represents the configuration for a readiness probe.
type ReadinessProbe struct {
	Exec Exec `yaml:"exec"`
}

// RegisterInstanceWrapper represents the registration information for a service.
type RegisterInstanceWrapper struct {
	Wrapper     string      // Wrapper is the name of the wrapper function.
	ServiceName interface{} // ServiceName is the name of the service.
	IP          interface{} // IP is the IP address of the service.
	Port        interface{} // Port is the port number of the service.
}

// Requests represents the resource requests for a container.
type Requests struct {
	Cpu    string `yaml:"cpu"`    // Cpu represents the CPU limit for the task.
	Memory string `yaml:"memory"` // Memory represents the memory limit for the task.
}

// Resources represents the resource requirements for a particular component.
type Resources struct {
	Requests Requests `yaml:"requests"` // Requests specifies the resource requests for the component.
	Limits   Limits   `yaml:"limits"`   // Limits specifies the resource limits for the component.
}

// SelectInstanceWrapper represents information about a selection.
type SelectInstanceWrapper struct {
	Wrapper     string      // Wrapper is the name of the wrapper.
	ServiceName interface{} // ServiceName is the name of the service.
}

// ServiceInfo represents information about a service.
type ServiceInfo struct {
	Application string // Application represents the name of the application.
	IP          string // IP represents the IP address of the service.
	Port        string // Port represents the port number of the service.
}

// Spec represents the specification of a resource.
type Spec struct {
	Template Template `yaml:"template"`
}

// TCPManifest represents the manifest for a TCP service.
type TCPManifest struct {
	Service  string       `json:"service"`  // Name of the service.
	Version  string       `json:"version"`  // Version of the service.
	Requests []TCPRequest `json:"requests"` // List of TCP requests.
}

// TCPRequest represents a TCP request.
type TCPRequest struct {
	Type string `json:"type"` // Type represents the type of the TCP request.
	URL  string `json:"url"`  // URL represents the URL of the TCP request.
	Name string `json:"name"` // Name represents the name of the TCP request.
	Port string `json:"port"` // Port represents the port number of the TCP request.
}

// Template represents a template object.
type Template struct {
	// Metadata contains the metadata of the template.
	Metadata TemplateMetadata `yaml:"metadata"`
	// Spec contains the specification of the template.
	Spec TemplateSpec `yaml:"spec"`
}

// TemplateMetadata represents the metadata of a template.
type TemplateMetadata struct {
	// Labels contains the labels associated with the template.
	Labels TemplateMetadataLabels `yaml:"labels"`
}

// TemplateMetadataLabels represents the metadata labels for a template.
type TemplateMetadataLabels struct {
	// App represents the application name.
	App string `yaml:"app"`
}

// TemplateSpec represents a template specification.
type TemplateSpec struct {
	Containers []Containers `yaml:"containers"`
}

// WrapperParams represents the parameters for a wrapper.
type WrapperParams struct {
	// Position represents the position at which the argument is passed into the wrapper.
	Position int
}

// Yaml2Go represents the structure of a YAML file converted to Go.
type Yaml2Go struct {
	// ApiVersion represents the API version of the YAML file.
	ApiVersion string `yaml:"apiVersion"`

	// Kind represents the kind of the YAML file.
	Kind string `yaml:"kind"`

	// Metadata represents the metadata of the YAML file.
	Metadata Metadata `yaml:"metadata"`

	// Spec represents the specification of the YAML file.
	Spec Spec `yaml:"spec"`
}
