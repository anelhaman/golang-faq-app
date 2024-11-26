package redis

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	once   sync.Once
)

// InitializeRedis initializes the Redis client as a singleton
func InitializeRedis() *redis.Client {
	once.Do(func() {
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")

		client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		})

		if err := client.Ping(client.Context()).Err(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}
		fmt.Println("Connected to Redis!")
	})
	return client
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	if client == nil {
		panic("Redis client is not initialized. Call InitializeRedis() first.")
	}
	return client
}
