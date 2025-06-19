package main

import (
	"fmt"
	"git.exclouds.org/ww/greet_ws/core"
	"git.exclouds.org/ww/greet_ws/global"
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
	//初始化
	core.RunWindowsServer()

	var GreetService = service.NewGreetService()

	// 初始化 Cron 调度器
	c := cron.New(cron.WithSeconds())

	// 设置每天 0 点执行定时任务
	_, err := c.AddFunc("0 * * * * *", func() {
		// 在每天的 0 点执行 GreetTaskV2
		global.LOGGER.Info(fmt.Sprintf("打招呼任务开始"))
		GreetService.GreetTaskV2()
	})

	if err != nil {
		fmt.Println("Error adding cron job:", err)
	}

	defer c.Stop() // 程序退出时停止 Cron
	// 启动 Cron 调度器
	c.Start()

	select {}

}
