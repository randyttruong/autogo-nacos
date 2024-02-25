package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// discoverHelloService uses Nacos to discover the Hello Service.
func discoverHelloService(namingClient naming_client.INamingClient) (string, error) {
	// Discover the Hello Service
	instances, err := namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: "HelloService",
		HealthyOnly: true,
	})
	if err != nil {
		return "", err
	}

	if len(instances) == 0 {
		return "", fmt.Errorf("no instances found for HelloService")
	}

	// Use the first available instance
	instance := instances[0]
	return fmt.Sprintf("http://%s", instance.Ip), nil
	// return fmt.Sprintf("http://%s:%d", instance.Ip, instance.Port), nil
}

func main() {
	log.Println("Starting Caller Service")
	// Nacos server configuration
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("host.docker.internal", 8848),
	}

	// Client configuration
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
	}

	// Create a naming client for service discovery
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Println("Failed to create naming client: %v", err)
	}

	// Discover the Hello Service URL
	helloServiceURL, err := discoverHelloService(namingClient)
	if err != nil {
		log.Println("Failed to discover Hello Service: %v", err)
	}

	// Get the username to greet (you can modify this part to get the username from different sources)
	username := os.Getenv("USERNAME")
	if username == "" {
		username = "DefaultUser"
	}

	// Call the Hello Service
	// 
	response, err := http.Get(fmt.Sprintf("%s/hello?username=%s", helloServiceURL, username))
	if err != nil {
		log.Println("Failed to call Hello Service: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to read response: %v", err)
	}

	fmt.Printf("Response from Hello Service: %s\n", string(body))
}
