package main

import (
	"os"

	"nacos-sdk-go-example/pkg/provider"
)

func main() {
	if err := provider.Execute(); err != nil {
		os.Exit(1)
	}
}
