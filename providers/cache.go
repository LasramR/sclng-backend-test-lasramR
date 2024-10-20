package providers

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Allow to perform cache related operations
type CacheProvider interface {
	// Retrieve and Unmarshall an element from the cache
	GetUnmarshalled(ctx context.Context, key string, unmarshalledPayload any) error
	// Marshal and Set an element in the cache with the given expiresIn duration
	SetMarshalled(ctx context.Context, key string, value any, expiresIn time.Duration) error
}

// IoC of the Redis client, we dont rely on HSet and struct tags because we don't want to be tighly coupled to redis
type RedisClient struct {
	Get func(context.Context, string) *redis.StringCmd
	Set func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type redisCacheProvider struct {
	client *RedisClient
}

func (r *redisCacheProvider) GetUnmarshalled(ctx context.Context, key string, unmarshalledPayload any) error {
	payload, err := r.client.Get(ctx, key).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(payload), unmarshalledPayload)
}

func (r *redisCacheProvider) SetMarshalled(ctx context.Context, key string, value any, expiresIn time.Duration) error {
	marshalled, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, marshalled, expiresIn).Err()
}

func NewRedisCacheProvider(redisClient *RedisClient) CacheProvider {
	return &redisCacheProvider{
		client: &RedisClient{
			Get: redisClient.Get,
			Set: redisClient.Set,
		},
	}
}
