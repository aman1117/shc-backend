package services

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO: study about this
// what is the meaning of context.Background()?✅
var ctx = context.Background()

type RedisService struct {
	client *redis.Client
}

func NewRedisService() *RedisService {
	// what does below line do?✅
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	return &RedisService{client: redis.NewClient(opts)}
}

func (rs *RedisService) GetCache(key string) (any, error) {
	// what does belowline do? -> removes leading and trailing whitespaces from the string✅
	key = strings.TrimSpace(key)

	if key == "" {
		return nil, errors.New("key is required")
	}
	// what does below line do? ✅
	val, err := rs.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (rs *RedisService) SetCache(key string, value interface{}, expires time.Duration) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("key is required")
	}

	// what does below line do?✅
	err := rs.client.Set(ctx, key, value, expires).Err()
	if err != nil {
		return err
	}
	return nil
}
