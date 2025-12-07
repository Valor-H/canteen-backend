package cache

import (
	"context"
	"canteen/internal/infrastructure/config"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisOnce sync.Once
	redisClient *redis.Client
)

// InitRedis 初始化Redis连接，从配置文件读取配置
func InitRedis() {
	redisOnce.Do(func() {
		redisConfig := &redis.Options{
			Addr:         config.GetString("redis.host") + ":" + config.GetString("redis.port"),
			Password:     config.GetString("redis.password"),
			DB:           config.GetInt("redis.db"),
			PoolSize:     config.GetInt("redis.pool_size"),
			MinIdleConns: config.GetInt("redis.min_idle_conns"),
		}
		redisClient = redis.NewClient(redisConfig)
		
		// 测试连接
		ctx := context.Background()
		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			log.Printf("Failed to connect to Redis: %v", err)
		} else {
			log.Println("Successfully connected to Redis")
		}
	})
}

// RedisClient 返回Redis客户端实例
func RedisClient() *redis.Client {
	if redisClient == nil {
		InitRedis()
	}
	return redisClient
}