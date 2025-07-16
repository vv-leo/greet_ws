package utils

/*
负载均衡使用示例：

1. 基本使用（默认轮询）：
   consulUtils := &ConsulUtils{}
   url, err := consulUtils.GetChatHistoryServiceURL()

2. 指定负载均衡策略：
   // 轮询模式 - 适合多节点均匀分布
   url, err := consulUtils.GetChatHistoryServiceURLWithStrategy(RoundRobin)

   // 随机模式 - 适合无状态服务
   url, err := consulUtils.GetChatHistoryServiceURLWithStrategy(Random)

3. 查看所有可用实例：
   instances, err := consulUtils.GetAllServiceInstances("chat-history-service")

多节点负载均衡效果：
- 如果有 3 个 ChatHistory 节点：192.168.1.10:8080, 192.168.1.11:8080, 192.168.1.12:8080
- 轮询模式会按顺序分配：第1次请求->节点1，第2次->节点2，第3次->节点3，第4次->节点1...
- 随机模式会随机选择任意一个健康节点
*/

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"git.exclouds.org/ww/greet_ws/global"
	"github.com/hashicorp/consul/api"
)

type ConsulUtils struct{}

// 负载均衡策略
type LoadBalanceStrategy int

const (
	RoundRobin     LoadBalanceStrategy = iota // 轮询
	Random                                    // 随机
	WeightedRandom                            // 加权随机（未来可扩展）
)

// 服务实例计数器，用于轮询负载均衡
var (
	serviceCounters = make(map[string]*ServiceCounter)
	counterMutex    sync.RWMutex
)

type ServiceCounter struct {
	count int64
	mutex sync.Mutex
}

// getServiceCounter 获取或创建服务计数器
func getServiceCounter(serviceName string) *ServiceCounter {
	counterMutex.RLock()
	counter, exists := serviceCounters[serviceName]
	counterMutex.RUnlock()

	if !exists {
		counterMutex.Lock()
		// 再次检查，避免重复创建
		if counter, exists = serviceCounters[serviceName]; !exists {
			counter = &ServiceCounter{count: 0}
			serviceCounters[serviceName] = counter
		}
		counterMutex.Unlock()
	}
	return counter
}

// GetServiceURL 动态获取指定服务的地址（默认轮询负载均衡）
func (c *ConsulUtils) GetServiceURL(serviceName string) (string, error) {
	return c.GetServiceURLWithStrategy(serviceName, RoundRobin)
}

// GetServiceURLWithStrategy 使用指定负载均衡策略获取服务地址
func (c *ConsulUtils) GetServiceURLWithStrategy(serviceName string, strategy LoadBalanceStrategy) (string, error) {
	// 创建 Consul 客户端
	consulClient, err := c.createConsulClient()
	if err != nil {
		return "", err
	}

	// 查询健康的服务实例
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("查询服务 %s 失败: %v", serviceName, err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("未找到可用的 %s 服务实例", serviceName)
	}

	// 根据策略选择服务实例
	selectedService := c.selectServiceByStrategy(services, serviceName, strategy)
	baseURL := fmt.Sprintf("http://%s:%d", selectedService.Service.Address, selectedService.Service.Port)

	global.LOGGER.Info(fmt.Sprintf("负载均衡选择服务 %s 地址: %s (策略: %v, 可用实例数: %d)",
		serviceName, baseURL, strategy, len(services)))
	return baseURL, nil
}

// createConsulClient 创建 Consul 客户端
func (c *ConsulUtils) createConsulClient() (*api.Client, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = global.SERVER_CONFIG.ConsulConfig.Address
	if global.SERVER_CONFIG.ConsulConfig.Datacenter != "" {
		consulConfig.Datacenter = global.SERVER_CONFIG.ConsulConfig.Datacenter
	}

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("创建 Consul 客户端失败: %v", err)
	}
	return consulClient, nil
}

// selectServiceByStrategy 根据负载均衡策略选择服务实例
func (c *ConsulUtils) selectServiceByStrategy(services []*api.ServiceEntry, serviceName string, strategy LoadBalanceStrategy) *api.ServiceEntry {
	switch strategy {
	case RoundRobin:
		return c.selectByRoundRobin(services, serviceName)
	case Random:
		return c.selectByRandom(services)
	default:
		// 默认使用轮询
		return c.selectByRoundRobin(services, serviceName)
	}
}

// selectByRoundRobin 轮询选择服务实例
func (c *ConsulUtils) selectByRoundRobin(services []*api.ServiceEntry, serviceName string) *api.ServiceEntry {
	counter := getServiceCounter(serviceName)
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	index := counter.count % int64(len(services))
	counter.count++

	global.LOGGER.Debug(fmt.Sprintf("轮询负载均衡: 服务 %s, 当前计数: %d, 选择索引: %d",
		serviceName, counter.count, index))

	return services[index]
}

// selectByRandom 随机选择服务实例
func (c *ConsulUtils) selectByRandom(services []*api.ServiceEntry) *api.ServiceEntry {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(services))

	global.LOGGER.Debug(fmt.Sprintf("随机负载均衡: 选择索引: %d", index))

	return services[index]
}

// GetChatHistoryServiceURLWithStrategy 使用指定策略获取 ChatHistory 服务地址
func (c *ConsulUtils) GetChatHistoryServiceURLWithStrategy(strategy LoadBalanceStrategy) (string, error) {
	return c.GetServiceURLWithStrategy("chat-history", strategy)
	//return "http://23.148.24.195:9991", nil
}

// GetChatSocketServiceURLWithStrategy 使用指定策略获取 ChatSocket 服务地址
func (c *ConsulUtils) GetChatSocketServiceURLWithStrategy(strategy LoadBalanceStrategy) (string, error) {
	return c.GetServiceURLWithStrategy("chat-socket", strategy)
	//return "http://23.148.24.195:7001", nil
}

// GetChatWorkServiceURLWithStrategy 使用指定策略获取 ChatWork 服务地址
func (c *ConsulUtils) GetChatWorkServiceURLWithStrategy(strategy LoadBalanceStrategy) (string, error) {
	return c.GetServiceURLWithStrategy("chat_work_ws", strategy)
	//return "http://23.148.24.195:11000", nil
}

// GetChatAccountServiceURLWithStrategy 使用指定策略获取 ChatAccount 服务地址
func (c *ConsulUtils) GetChatAccountServiceURLWithStrategy(strategy LoadBalanceStrategy) (string, error) {
	return c.GetServiceURLWithStrategy("chat_accountPool", strategy)
	//return "http://23.148.24.195:10100", nil
}

// GetAllServiceInstances 获取所有健康的服务实例信息（用于监控和调试）
func (c *ConsulUtils) GetAllServiceInstances(serviceName string) ([]*api.ServiceEntry, error) {
	consulClient, err := c.createConsulClient()
	if err != nil {
		return nil, err
	}

	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("查询服务 %s 失败: %v", serviceName, err)
	}

	global.LOGGER.Info(fmt.Sprintf("服务 %s 当前有 %d 个健康实例", serviceName, len(services)))
	for i, service := range services {
		global.LOGGER.Info(fmt.Sprintf("实例 %d: %s:%d (ID: %s)",
			i, service.Service.Address, service.Service.Port, service.Service.ID))
	}

	return services, nil
}
