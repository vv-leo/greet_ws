package main

import (
	"fmt"
	"time"

	"git.exclouds.org/ww/greet_ws/core"
	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/initialize"
	"git.exclouds.org/ww/greet_ws/service"
	"github.com/robfig/cron/v3"
)

// @title global项目
// @version 0.1
// @description json访问：<a href="/docs/swagger.json" target="_blank">swagger.json</a>
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	// 1. 初始化外部接口
	global.InitExternalApiUrl()
	global.LOGGER.Info("外部接口初始化完成")

	// 2. 初始化Redis
	if global.SERVER_CONFIG.SystemConfig.UseMultipoint {
		initialize.RedisInitialize()
		global.LOGGER.Info("Redis初始化完成")
	}

	// 3. 启动服务（HTTP服务器+Consul注册）
	go func() {
		core.RunWindowsServer()
	}()
	time.Sleep(3 * time.Second) // 等待服务启动
	global.LOGGER.Info("HTTP服务启动完成")

	// 4. 启动定时任务
	startScheduler()

	// 保持程序运行
	select {}
}

// startScheduler 启动定时任务
func startScheduler() {
	global.LOGGER.Info("启动定时任务...")

	greetService := service.NewGreetService()
	scheduler := cron.New(cron.WithSeconds())

	// 每分钟执行一次打招呼任务
	_, err := scheduler.AddFunc("0 * * * * *", func() {
		global.LOGGER.Info("开始执行打招呼任务")
		start := time.Now()
		greetService.GreetTask()
		duration := time.Since(start)
		global.LOGGER.Info(fmt.Sprintf("打招呼任务执行完成，耗时: %v", duration))
	})

	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("定时任务启动失败: %v", err))
		return
	}

	scheduler.Start()
	global.LOGGER.Info("定时任务启动成功")
}
