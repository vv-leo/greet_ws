package utils

import (
	"fmt"
	"math/rand"
	"time"

	"git.exclouds.org/ww/greet_ws/global"
	"github.com/hashicorp/consul/api"
)

type ServiceDiscoveryUtils struct {
	client *api.Client
}

func NewServiceDiscoveryUtils() (*ServiceDiscoveryUtils, error) {
	config := api.DefaultConfig()
	config.Address = global.SERVER_CONFIG.ConsulConfig.Address

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建Consul客户端失败: %v", err)
	}

	return &ServiceDiscoveryUtils{client: client}, nil
}

// GetServiceURL 获取服务URL（负载均衡）
func (s *ServiceDiscoveryUtils) GetServiceURL(serviceName string) (string, error) {
	// 获取健康的服务实例
	services, _, err := s.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("获取服务实例失败: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("没有找到健康的服务实例: %s", serviceName)
	}

	// 随机选择一个实例（简单的负载均衡）
	rand.Seed(time.Now().UnixNano())
	selectedService := services[rand.Intn(len(services))]

	serviceURL := fmt.Sprintf("http://%s:%d",
		selectedService.Service.Address,
		selectedService.Service.Port)

	return serviceURL, nil
}

// GetAllServiceInstances 获取所有服务实例
func (s *ServiceDiscoveryUtils) GetAllServiceInstances(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := s.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("获取服务实例失败: %v", err)
	}
	return services, nil
}
