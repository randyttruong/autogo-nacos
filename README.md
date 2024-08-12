# Nacos Golang Static Analyser

This is a utility written in Go and Python for automatically generating microsegmentation policies for applications that utilize Kubernetes container orchestration and Nacos Service Mesh. It analyses an application's source code, generates a `TCPManifest` containing the name of the service, the version of the service, the TCP calls it makes, and finally generates Kubernetes security policies by generating a service graph of the application topology.

## Prerequisites

- The repository to be analyzed must be a valid Kubernetes project with a YAML config file.
- Golang 1.22.0

## Instructions

1. Create an `input` folder then place the Kubernetes project in the `input` folder.
2. Create an `output` folder for maniect file to be generated into
3. Edit the `main.go` function's output prefix to your preferred name.
4. Navigate to the `static_analyser` directory.
  ```
  cd static_analyser
  ```
4. Build the project.
  ```
  go build -o bin/static_analyser ./cmd/static_analyser
  ```
  For Windows, run this instead:
  ```
  go build -o bin/static_analyser.exe ./cmd/static_analyser
  ```
5. Run the static analyser.
  ```
  ./bin/static_analyser
  ```
  For Windows, run this instead:
  ```
  ./bin/static_analyser.exe
  ```
6. The output will be placed in the `output` folder with the output prefix you specified in step 2.

## Output

The output is a TCPManifest containing the name of the service, the version of the service, and the TCP calls it made.
