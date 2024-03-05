package provider

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"nacos-sdk-go-example/pkg/config"
	"nacos-sdk-go-example/pkg/util/name"
	"nacos-sdk-go-example/pkg/util/naming"
	"nacos-sdk-go-example/pkg/util/random"
)

const (
	defaultNamingClientTime   = 1 * time.Second
	defaultProviderConfigPath = "/configs/provider"
)

func Execute() error {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return err
	}
	configPath := basePath + defaultProviderConfigPath
	config.ReadConfig("config", configPath, "yaml")
	nameAddr := config.ConfigMessage.Basic.NameServerAddr
	providerName, err := name.ReadName(nameAddr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// get serverConfigs & clientConfig
	serverConfigs := naming.GetServerConfigs(config.ConfigMessage.Server.IpAddr, config.ConfigMessage.Server.Port)
	clientConfig := naming.GetClientConfig()
	// new naming client
	client, err := naming.NewNamingClient(serverConfigs, clientConfig)
	for err != nil {
		time.Sleep(defaultNamingClientTime)
		client, err = naming.NewNamingClient(serverConfigs, clientConfig)
	}
	scope := config.ConfigMessage.Client.Scope

	// register instance
	registerInstanceParam := vo.RegisterInstanceParam{
		Ip:          config.ConfigMessage.Basic.InstanceIp + strconv.Itoa(random.Numb(scope)),
		Port:        config.ConfigMessage.Basic.InstancePort + uint64(random.Numb(scope)),
		ServiceName: providerName,
		Weight:      10,
		ClusterName: config.ConfigMessage.Basic.InstanceClusterName,
		GroupName:   "group-A",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	err = naming.RegisterServiceInstance(client, registerInstanceParam)
	for {
		if err != nil {
			err = naming.RegisterServiceInstance(client, registerInstanceParam)
		} else {
			break
		}
	}

	// sleep
	time.Sleep(360000 * time.Second)
	return nil
}
