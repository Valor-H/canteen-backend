package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

type dailyLogWriter struct {
	mu          sync.Mutex
	currentDate string
	writer      io.Writer
}

func (w *dailyLogWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != w.currentDate {
		logFilePath := fmt.Sprintf("./logs/%s.log", today)
		file := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    10,
			MaxBackups: 7,
			MaxAge:     60,
			Compress:   true,
		}
		w.writer = io.MultiWriter(os.Stdout, file)
		w.currentDate = today
	}
	return w.writer.Write(p)
}

type dailyIllegalLogWriter struct {
	mu          sync.Mutex
	currentDate string
	writer      io.Writer
}

func (w *dailyIllegalLogWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != w.currentDate {
		logPath := fmt.Sprintf("./logs/illegal-%s.log", today)
		file := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    10,
			MaxBackups: 7,
			MaxAge:     60,
			Compress:   true,
		}
		w.writer = file
		w.currentDate = today
	}
	return w.writer.Write(p)
}

var illegalLogger *log.Logger

func init() {
	// 初始化非法请求日志写入器
	dlw := &dailyIllegalLogWriter{
		currentDate: time.Now().Format("2006-01-02"),
		writer: &lumberjack.Logger{
			Filename:   fmt.Sprintf("./logs/illegal-%s.log", time.Now().Format("2006-01-02")),
			MaxSize:    10,
			MaxBackups: 7,
			MaxAge:     60,
			Compress:   true,
		},
	}
	illegalLogger = log.New(dlw, "", log.LstdFlags)
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	// 初始化主日志动态写入器
	dlw := &dailyLogWriter{
		currentDate: time.Now().Format("2006-01-02"),
		writer: io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   fmt.Sprintf("./logs/%s.log", time.Now().Format("2006-01-02")),
			MaxSize:    10,
			MaxBackups: 7,
			MaxAge:     60,
			Compress:   true,
		}),
	}

	// 设置全局日志输出
	log.SetOutput(dlw)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		log.Printf("[%-d]\t%-v\t[%-s]\t[%-s]\t[%-s]\tUser-Agent: %-s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			userAgent,
		)
	}
}

// GetIllegalLogger 获取非法请求日志记录器
func GetIllegalLogger() *log.Logger {
	return illegalLogger
}