package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"canteen/internal/infrastructure/config"
	"canteen/internal/router"
	"canteen/internal/app"
	"canteen/pkg/utils"
)

func main() {
	// 验证许可证
	if !utils.ValidateLicense() {
		log.Fatal("License validation failed. Exiting...")
	}
	
	// 初始化配置
	config.InitConfig()
	
	// 创建并初始化应用
	app := app.NewApplication()
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	
	// 确保在程序退出时关闭应用
	defer func() {
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during application shutdown: %v", err)
		}
	}()
	
	// 启动后台任务
	app.StartBackgroundTasks()
	
	// 启动HTTP服务器
	port := config.GetString("server.port")
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router.Config(),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	// 在goroutine中启动服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server failed to start: %v", err)
		}
	}()
	
	// 等待中断信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	
	log.Println("Shutting down the server...")
	
	// 创建一个带有超时的上下文，用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 优雅关闭HTTP服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited")
}