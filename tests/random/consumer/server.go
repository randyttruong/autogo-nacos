package consumer

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"nacos-sdk-go-example/pkg/config"
	"nacos-sdk-go-example/pkg/util/name"
	"nacos-sdk-go-example/pkg/util/naming"
	"nacos-sdk-go-example/pkg/util/random"
)

var address string

const (
	defaultNamingClientTime   = 1 * time.Second
	defaultConsumerConfigPath = "/configs/consumer"
)

func Execute() error {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return err
	}
	configPath := basePath + defaultConsumerConfigPath
	config.ReadConfig("config", configPath, "yaml")
	nameAddr := config.ConfigMessage.Basic.NameServerAddr
	consumerName, err := name.ReadName(nameAddr)
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
		ServiceName: consumerName,
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

	// subscribe
	basicServiceName := config.ConfigMessage.Client.ServiceName
	scope = config.ConfigMessage.Client.Scope
	subscribeNum := config.ConfigMessage.Basic.SubscribeInstanceCount
	for i := 1; i <= subscribeNum; i++ {
		subscribeName := random.ServiceName(basicServiceName, scope)
		fmt.Println(subscribeName)
		err = naming.Subscribe(client, &vo.SubscribeParam{
			ServiceName: subscribeName,
			Clusters:    []string{config.ConfigMessage.Basic.InstanceClusterName},
			GroupName:   "group-A",
			SubscribeCallback: func(instances []model.Instance, err error) {
				if len(instances) != 0 {
					address = fmt.Sprintf("http://%s:%v/", instances[0].Ip, instances[0].Port)
				} else {
					address = ""
				}
				fmt.Println(fmt.Sprintf("callback service: %+v return instance:%+v \n", "provider", instances))
			},
		})
		if err != nil {
			panic(err)
		}
	}
	// sleep
	time.Sleep(360000 * time.Second)
	return nil
}
