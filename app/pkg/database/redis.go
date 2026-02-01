package database

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisMu     sync.RWMutex
)

func ConnectRedis(ctx context.Context, addr string) error {
	redisMu.Lock()
	defer redisMu.Unlock()

	redisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	return redisClient.Ping(ctx).Err()
}

func Redis() *redis.Client {
	redisMu.RLock()
	defer redisMu.RUnlock()

	return redisClient
}

func CloseRedis() error {
	redisMu.Lock()
	defer redisMu.Unlock()

	if redisClient != nil {
		return redisClient.Close()
	}

	return nil
}

func StartRedisHealthCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
            case <-ctx.Done():
                return

            case <-ticker.C:
                if err := Redis().Ping(ctx).Err(); err != nil {
                    slog.Error("redis health check failed", "err", err)
                }
		}
	}
}