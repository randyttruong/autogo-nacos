# Nacos Golang Static Analyser

This project is a static analyser for Kubernetes projects using `nacos-sdk-go` and written in Go. It analyses the repository and generates a TCPManifest containing the name of the service, the version of the service, and the TCP calls it made.

## Prerequisites

- The repository to be analyzed must be a valid Kubernetes project with a YAML config file.
- Golang 1.22.0

## Instructions

1. Place the Kubernetes project in the `input` folder.
2. Edit the `main.go` function's output prefix to your preferred name.
3. Navigate to the `static_analyser` directory.
  ```
  cd static_analyser
  ```
4. Build the project.
  ```
  go build -o bin/static_analyser ./cmd/static_analyser
  ```
5. Run the static analyser.
  ```
  ./bin/static_analyser
  ```
6. The output will be placed in the `output` folder with the output prefix you specified in step 2.

## Output

The output is a TCPManifest containing the name of the service, the version of the service, and the TCP calls it made.
