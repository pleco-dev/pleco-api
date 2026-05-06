package appsetup

import (
	"context"
	"fmt"
	"log"
	"strings"

	"pleco-api/internal/cache"
	"pleco-api/internal/config"

	"github.com/redis/go-redis/v9"
)

func newCacheStore(cfg config.AppConfig) cache.Store {
	redisURL := redisConnectionURL(cfg)
	if redisURL == "" {
		log.Println("cache: using in-memory store")
		return cache.NewMemoryStore()
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("cache: invalid Redis config, using in-memory store: %v", err)
		return cache.NewMemoryStore()
	}

	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("cache: Redis unavailable (%v), using in-memory store", err)
		_ = rdb.Close()
		return cache.NewMemoryStore()
	}

	log.Println("cache: using Redis store")
	return cache.NewRedisStore(rdb)
}

func redisConnectionURL(cfg config.AppConfig) string {
	if strings.TrimSpace(cfg.RedisURL) != "" {
		return strings.TrimSpace(cfg.RedisURL)
	}

	host := strings.TrimSpace(config.GetEnv("REDIS_HOST", ""))
	if host == "" {
		return ""
	}
	port := strings.TrimSpace(config.GetEnv("REDIS_PORT", "6379"))
	password := strings.TrimSpace(config.GetEnv("REDIS_PASSWORD", ""))
	db := strings.TrimSpace(config.GetEnv("REDIS_DB", "0"))
	if password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/%s", password, host, port, db)
	}
	return fmt.Sprintf("redis://%s:%s/%s", host, port, db)
}
