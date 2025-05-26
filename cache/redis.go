package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
)

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. Caching will be disabled.", err)
		Client = nil
	} else {
		log.Println("Successfully connected to Redis")
	}
}

func Set(key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return nil
	}
	ctx := context.Background()
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, json, expiration).Err()
}

func Get(key string, dest interface{}) error {
	if Client == nil {
		return redis.Nil
	}
	ctx := context.Background()
	val, err := Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func Delete(key string) error {
	if Client == nil {
		return nil
	}
	ctx := context.Background()
	return Client.Del(ctx, key).Err()
}

func Exists(key string) (bool, error) {
	if Client == nil {
		return false, nil
	}
	ctx := context.Background()
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}
