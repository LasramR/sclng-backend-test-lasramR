package providers

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func MockRedisCacheClient(val string, err error, fakeCache map[string]interface{}) *RedisClient {
	return &RedisClient{
		Get: func(ctx context.Context, s string) *redis.StringCmd {
			cmd := &redis.StringCmd{}
			if err != nil {
				cmd.SetErr(err)
			} else {
				cmd.SetVal(val)
			}
			return cmd
		},
		Set: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			cmd := &redis.StatusCmd{}
			if err != nil {
				cmd.SetErr(err)
			} else if fakeCache != nil {
				fakeCache[key] = value
			}
			return cmd
		},
	}
}

type TestStruct struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func TestRedisGetUnmarshelled_WithCachedValue(t *testing.T) {
	cacheProvider := NewRedisCacheProvider(MockRedisCacheClient(`{"key":"somefield","value":1}`, nil, nil))

	ctx := context.Background()
	var s TestStruct
	err := cacheProvider.GetUnmarshalled(ctx, "some key", &s)

	if err != nil {
		t.Fatalf("When redis cache returns a valid marshalled struct no error should be thrown")
	}

	expected := TestStruct{
		Key:   "somefield",
		Value: 1,
	}

	if !reflect.DeepEqual(s, expected) {
		t.Fatalf("Unmarshalled value from redis cache should equals to expected")
	}
}

func TestRedisGetUnmarshelled_WithoutCachedValue(t *testing.T) {
	cacheProvider := NewRedisCacheProvider(MockRedisCacheClient(``, errors.New("not cached"), nil))

	ctx := context.Background()
	var s TestStruct
	err := cacheProvider.GetUnmarshalled(ctx, "some key", &s)

	if err == nil {
		t.Fatalf("When redis cache returns doesnt returns valid marshalled struct it should return an error")
	}
}

func TestRedisSetMarshalled(t *testing.T) {
	fakeCache := make(map[string]interface{})
	cacheProvider := NewRedisCacheProvider(MockRedisCacheClient(`some key`, nil, fakeCache))

	ctx := context.Background()
	s := TestStruct{
		Key:   "somefield",
		Value: 1,
	}

	err := cacheProvider.SetMarshalled(ctx, "some key", s, time.Hour)

	if err != nil {
		t.Fatalf("Setting up a valid marshallable value in cache should not return an error")
	}

	expected := `{"key":"somefield","value":1}`
	if !reflect.DeepEqual(fakeCache["some key"], []byte(expected)) {
		t.Fatalf("Cached value should be equal to expected")
	}
}
