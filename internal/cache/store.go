package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	GetJSON(ctx context.Context, key string, dst any) (bool, error)
	SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeletePrefix(ctx context.Context, prefix string) error
}

type NoopStore struct{}

func (NoopStore) GetJSON(context.Context, string, any) (bool, error) { return false, nil }
func (NoopStore) SetJSON(context.Context, string, any, time.Duration) error {
	return nil
}
func (NoopStore) Delete(context.Context, ...string) error    { return nil }
func (NoopStore) DeletePrefix(context.Context, string) error { return nil }

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (s *RedisStore) GetJSON(ctx context.Context, key string, dst any) (bool, error) {
	raw, err := s.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(raw, dst)
}

func (s *RedisStore) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, raw, ttl).Err()
}

func (s *RedisStore) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, keys...).Err()
}

func (s *RedisStore) DeletePrefix(ctx context.Context, prefix string) error {
	iter := s.client.Scan(ctx, 0, prefix+"*", 100).Iterator()
	keys := make([]string, 0, 100)
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
		if len(keys) >= 100 {
			if err := s.Delete(ctx, keys...); err != nil {
				return err
			}
			keys = keys[:0]
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return s.Delete(ctx, keys...)
}

type memoryEntry struct {
	value     []byte
	expiresAt time.Time
}

type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]memoryEntry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: make(map[string]memoryEntry)}
}

func (s *MemoryStore) GetJSON(_ context.Context, key string, dst any) (bool, error) {
	s.mu.RLock()
	entry, ok := s.items[key]
	s.mu.RUnlock()
	if !ok {
		return false, nil
	}
	if time.Now().After(entry.expiresAt) {
		_ = s.Delete(context.Background(), key)
		return false, nil
	}
	return true, json.Unmarshal(entry.value, dst)
}

func (s *MemoryStore) SetJSON(_ context.Context, key string, value any, ttl time.Duration) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.items[key] = memoryEntry{value: raw, expiresAt: time.Now().Add(ttl)}
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) Delete(_ context.Context, keys ...string) error {
	s.mu.Lock()
	for _, key := range keys {
		delete(s.items, key)
	}
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) DeletePrefix(_ context.Context, prefix string) error {
	s.mu.Lock()
	for key := range s.items {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(s.items, key)
		}
	}
	s.mu.Unlock()
	return nil
}
