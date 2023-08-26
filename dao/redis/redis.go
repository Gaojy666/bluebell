package redis

import (
	"bluebell/settings"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize, // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	return err
}

func Close() {
	_ = client.Close()
}
