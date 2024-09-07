package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	once   sync.Once
	ctx    = context.Background()
)

// initialize 初始化 Redis 客户端
func initialize() {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "likkim2024",
			DB:       0,
		})

		// 测试连接
		_, err := client.Ping(ctx).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Redis client initialized")
	})
}

// GetClient 获取 Redis 客户端
func GetClient() *redis.Client {
	if client == nil {
		initialize()
	}
	return client
}

func DefaultCtx() context.Context {
	return context.Background()
}
