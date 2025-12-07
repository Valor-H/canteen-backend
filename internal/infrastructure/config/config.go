package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance *viper.Viper
)

// InitConfig 初始化配置
func InitConfig() {
	once.Do(func() {
		v := viper.New()

		// 设置配置文件的名称和类型
		v.SetConfigName("config") // 配置文件名（不带扩展名）
		v.SetConfigType("yaml")   // 配置文件类型

		// 添加配置文件路径
		exePath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get executable path: %v", err)
		}

		configDir := filepath.Join(exePath, "config") // 拼接 config 目录路径
		v.AddConfigPath(configDir)                    // 添加 config 目录作为配置文件路径

		// 读取配置文件
		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}

		// 将实例保存到全局变量
		instance = v
	})
}

// GetString 获取字符串类型的配置值
func GetString(key string) string {
	return instance.GetString(key)
}

// GetInt 获取整数类型的配置值
func GetInt(key string) int {
	return instance.GetInt(key)
}

// GetBool 获取布尔类型的配置值
func GetBool(key string) bool {
	return instance.GetBool(key)
}