package redis

import (
	rV8 "github.com/go-redis/redis/v8"
)

type redisBase struct {
	redisClient rV8.UniversalClient
}
