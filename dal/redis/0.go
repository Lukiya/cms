package redis

import (
	rV9 "github.com/redis/go-redis/v9"
)

type redisBase struct {
	redisClient rV9.UniversalClient
	redisKey    string
}
