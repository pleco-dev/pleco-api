package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimitStore implements a fixed-window counter per key using Redis INCR + TTL.
type RedisRateLimitStore struct {
	rdb    *redis.Client
	prefix string
}

var rateLimitAllowScript = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("PTTL", KEYS[1])
return {count, ttl}
`)

func NewRedisRateLimitStore(rdb *redis.Client) *RedisRateLimitStore {
	if rdb == nil {
		return nil
	}
	return &RedisRateLimitStore{rdb: rdb, prefix: "ratelimit:"}
}

func (s *RedisRateLimitStore) redisKey(key string) string {
	return s.prefix + key
}

func (s *RedisRateLimitStore) Allow(key string, limit int, window time.Duration, now time.Time) (bool, time.Time) {
	ctx := context.Background()
	rkey := s.redisKey(key)

	result, err := rateLimitAllowScript.Run(ctx, s.rdb, []string{rkey}, window.Milliseconds()).Result()
	if err != nil {
		// Fail open so a Redis outage does not hard-block the API.
		return true, now.Add(window)
	}

	values, ok := result.([]any)
	if !ok || len(values) != 2 {
		return true, now.Add(window)
	}

	count, err := redisInt64(values[0])
	if err != nil {
		return true, now.Add(window)
	}

	expiresAt := now.Add(window)
	ttlMillis, err := redisInt64(values[1])
	if err == nil && ttlMillis > 0 {
		expiresAt = now.Add(time.Duration(ttlMillis) * time.Millisecond)
	}

	return int(count) <= limit, expiresAt
}

func (s *RedisRateLimitStore) Close() error {
	if s == nil || s.rdb == nil {
		return nil
	}
	return s.rdb.Close()
}

func redisInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unexpected redis numeric type %T", value)
	}
}
