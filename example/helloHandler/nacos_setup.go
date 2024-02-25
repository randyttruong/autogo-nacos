package main

import (
    // "log"

    "github.com/nacos-group/nacos-sdk-go/clients"
    "github.com/nacos-group/nacos-sdk-go/common/constant"
    "github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosConfig struct {
	ServerIP   string
	ServerPort uint64
	Namespace  string
}

func RegisterService(nacosConfig NacosConfig, serviceName string, ip string, port uint64) error {
	// Server configurations
	sc := []constant.ServerConfig{
			*constant.NewServerConfig(nacosConfig.ServerIP, nacosConfig.ServerPort),
	}

	// Client configuration
	cc := constant.ClientConfig{
			NamespaceId:         nacosConfig.Namespace,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
	}

	// Create a naming client for service registration
	namingClient, err := clients.NewNamingClient(
			vo.NacosClientParam{
					ClientConfig:  &cc,
					ServerConfigs: sc,
			},
	)
	if err != nil {
			return err
	}

	// Register the service
	_, err = namingClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          ip,
			Port:        port,
			ServiceName: serviceName,
			Weight:      10,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
	})
	if err != nil {
			return err
	}

	return nil
}
