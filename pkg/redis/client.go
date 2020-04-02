package redis

import (
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

// GetClient get redis client.
func GetClient() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}

	return redisClient
}
