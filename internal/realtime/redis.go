package realtime

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func InitRedis(addr string) error {
	Client = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return Client.Ping(context.Background()).Err()
} 
