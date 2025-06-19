package initialize

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	cfg "git.exclouds.org/ww/greet_ws/config"
	"git.exclouds.org/ww/greet_ws/global"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var (
	consulClient *api.Client
	serviceID    string
	consulMutex  sync.Mutex
	isRegistered bool
)

// ConsulInitialize 初始化 Consul 客户端并注册服务
func ConsulInitialize(port int) error {
	consulMutex.Lock()
	defer consulMutex.Unlock()

	// 避免重复注册
	if isRegistered {
		return nil
	}

	// 检查 Consul 配置
	if global.SERVER_CONFIG.ConsulConfig.Address == "" {
		return fmt.Errorf("consul address not configured")
	}

	// 创建 Consul 客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = global.SERVER_CONFIG.ConsulConfig.Address
	if global.SERVER_CONFIG.ConsulConfig.Datacenter != "" {
		consulConfig.Datacenter = global.SERVER_CONFIG.ConsulConfig.Datacenter
	}

	var err error
	consulClient, err = api.NewClient(consulConfig)
	if err != nil {
		global.LOGGER.Error("Consul client creation failed", zap.Error(err))
		return err
	}

	// 获取服务信息
	serviceInfo := global.SERVER_CONFIG.ConsulConfig.Service

	// 生成服务ID
	hostname, _ := os.Hostname()
	serviceID = generateServiceID(serviceInfo, hostname, port)

	// 获取服务地址
	serviceAddress := getServiceAddress(serviceInfo.Address, hostname)

	// 构建健康检查配置
	checkURL := buildHealthCheckURL(serviceAddress, port, serviceInfo.Check.HTTP)
	checkConfig := buildHealthCheckConfig(serviceInfo.Check)
	checkConfig.HTTP = checkURL

	// 注册服务
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceInfo.Name,
		Tags:    serviceInfo.Tags,
		Port:    port,
		Address: serviceAddress,
		Check:   checkConfig,
	}

	if err := consulClient.Agent().ServiceRegister(registration); err != nil {
		global.LOGGER.Error("Service registration failed", zap.Error(err))
		return err
	}

	isRegistered = true
	global.LOGGER.Info("Service registered successfully",
		zap.String("service_id", serviceID),
		zap.String("service_name", serviceInfo.Name),
		zap.String("service_address", serviceAddress),
		zap.Int("port", port),
		zap.String("check_url", checkURL),
	)

	return nil
}

// DeregisterService 注销服务
func DeregisterService() error {
	consulMutex.Lock()
	defer consulMutex.Unlock()

	if !isRegistered || consulClient == nil || serviceID == "" {
		return nil
	}

	// 重试注销，确保服务被正确注销
	maxRetries := 3
	var err error

	for i := 0; i < maxRetries; i++ {
		err = consulClient.Agent().ServiceDeregister(serviceID)
		if err == nil {
			isRegistered = false
			global.LOGGER.Info("Service deregistered successfully", zap.String("service_id", serviceID))
			return nil
		}

		global.LOGGER.Warn("Service deregistration retry",
			zap.Int("attempt", i+1),
			zap.Error(err))

		time.Sleep(time.Second)
	}

	global.LOGGER.Error("Service deregistration failed after retries", zap.Error(err))
	return err
}

// GetPortFromEnv 从环境变量获取端口号
func GetPortFromEnv() int {
	portEnv := os.Getenv("SERVICE_PORT")
	if portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err == nil && port > 0 && port < 65536 {
			return port
		}
		global.LOGGER.Warn("Invalid port in environment variable, using default",
			zap.String("env_port", portEnv),
			zap.Error(err),
		)
	}

	return global.SERVER_CONFIG.SystemConfig.Addr
}

// 辅助函数：生成服务ID
func generateServiceID(serviceInfo cfg.ServiceInfo, hostname string, port int) string {
	// 使用配置中的ID模板或生成新ID
	podIP := os.Getenv("POD_IP")
	instanceID := os.Getenv("INSTANCE_ID")

	if serviceInfo.ID == "" {
		if instanceID != "" {
			return fmt.Sprintf("%s-%s", serviceInfo.Name, instanceID)
		} else if podIP != "" {
			return fmt.Sprintf("%s-%s", serviceInfo.Name, strings.ReplaceAll(podIP, ".", "-"))
		} else {
			return fmt.Sprintf("%s-%s-%d", serviceInfo.Name, hostname, port)
		}
	}

	// 变量替换
	id := serviceInfo.ID
	id = strings.ReplaceAll(id, "{hostname}", hostname)
	id = strings.ReplaceAll(id, "{port}", strconv.Itoa(port))
	if podIP != "" {
		id = strings.ReplaceAll(id, "{podIP}", podIP)
	}
	if instanceID != "" {
		id = strings.ReplaceAll(id, "{instanceID}", instanceID)
	}

	return id
}

// 辅助函数：获取服务地址
func getServiceAddress(configAddress string, hostname string) string {
	if configAddress != "" {
		return configAddress
	}

	// 优先使用环境变量中的IP
	podIP := os.Getenv("POD_IP")
	if podIP != "" {
		return podIP
	}

	// 尝试获取本机IP
	ip, err := getLocalIP()
	if err == nil {
		return ip
	}

	// 默认使用主机名
	return hostname
}

// 辅助函数：获取本机IP
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no suitable IP address found")
}

// 辅助函数：构建健康检查URL
func buildHealthCheckURL(serviceAddress string, port int, configPath string) string {
	if configPath == "" {
		return fmt.Sprintf("http://%s:%d/health", serviceAddress, port)
	}

	if configPath[0] == '/' {
		return fmt.Sprintf("http://%s:%d%s", serviceAddress, port, configPath)
	}

	return configPath
}

// 辅助函数：构建健康检查配置
func buildHealthCheckConfig(checkConfig cfg.ServiceCheck) *api.AgentServiceCheck {
	// 设置默认值
	interval := "10s"
	if checkConfig.Interval != "" {
		interval = checkConfig.Interval
	}

	timeout := "5s"
	if checkConfig.Timeout != "" {
		timeout = checkConfig.Timeout
	}

	deregister := "30s"
	if checkConfig.DeregisterCriticalServiceAfter != "" {
		deregister = checkConfig.DeregisterCriticalServiceAfter
	}

	// 构建健康检查URL
	checkURL := checkConfig.HTTP

	return &api.AgentServiceCheck{
		HTTP:                           checkURL,
		Interval:                       interval,
		Timeout:                        timeout,
		DeregisterCriticalServiceAfter: deregister,
	}
}
