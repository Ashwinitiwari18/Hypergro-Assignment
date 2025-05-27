package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Printf("Warning: REDIS_URL environment variable not set. Redis caching will be disabled.")
		return
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Warning: Error parsing Redis URL: %v. Redis caching will be disabled.", err)
		return
	}

	Client = redis.NewClient(opt)

	// Test connection
	_, err = Client.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. Caching will be disabled.", err)
		Client = nil
		return
	}

	log.Printf("Successfully connected to Redis!")
}

func Set(key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return Client.Set(Ctx, key, value, expiration).Err()
}

func Get(key string, dest interface{}) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return Client.Get(Ctx, key).Scan(dest)
}

func Delete(key string) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return Client.Del(Ctx, key).Err()
}

func Exists(key string) (bool, error) {
	if Client == nil {
		return false, nil
	}
	ctx := context.Background()
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}
