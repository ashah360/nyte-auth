package cache

import "github.com/go-redis/redis/v8"

func NewClient(addr, password string, DB int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       DB,
	})
}
