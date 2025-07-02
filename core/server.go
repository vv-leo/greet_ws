package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/initialize"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func RunWindowsServer() {
	// 初始化路由
	Router := initialize.RoutersInitialize()

	// 获取端口号（优先从环境变量读取）
	port := initialize.GetPortFromEnv()
	address := fmt.Sprintf(":%d", port)

	// 初始化服务器
	s := initServer(address, Router)

	// 设置HTTP服务器优雅关闭处理
	idleConnsClosed := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		global.LOGGER.Info("HTTP服务器接收到关闭信号，开始优雅关闭...")

		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// 优雅关闭HTTP服务器
		if err := s.Shutdown(ctx); err != nil {
			global.LOGGER.Error("HTTP服务器关闭失败", zap.Error(err))
		} else {
			global.LOGGER.Info("HTTP服务器已优雅关闭")
		}

		close(idleConnsClosed)
	}()

	// 注册到 Consul（如果配置了 Consul）
	if global.SERVER_CONFIG.ConsulConfig.Address != "" {
		err := initialize.ConsulInitialize(port)
		if err != nil {
			global.LOGGER.Warn("Consul 服务注册失败", zap.Error(err))
		} else {
			global.LOGGER.Info("Consul 服务注册成功")
		}
	}

	// 保证文本顺序输出
	time.Sleep(10 * time.Microsecond)
	global.LOGGER.Info("服务运行中", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 云控系统[ws-greet]
	当前版本:V0.1
	服务地址:%s
`, address)

	// 启动服务
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		global.LOGGER.Error("服务意外停止", zap.Error(err))
	}

	// 等待优雅关闭完成
	<-idleConnsClosed
}
