package main

import (
	"os"

	"nacos-sdk-go-example/pkg/consumer"
)

func main() {
	if err := consumer.Execute(); err != nil {
		os.Exit(1)
	}
}
