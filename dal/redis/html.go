package redis

import (
	"context"

	"github.com/Lukiya/cms/dal"
	rV8 "github.com/go-redis/redis/v8"
	"github.com/syncfuture/go/serr"
	"github.com/syncfuture/go/sredis"
)

func NewRedisHtmlCache(key string, redisConfig *sredis.RedisConfig) dal.IHtmlCacheDAL {
	r := new(RedisHtmlCache)
	r.redisClient = sredis.NewClient(redisConfig)
	r.key = key
	return r
}

type RedisHtmlCache struct {
	redisBase
}

func (x *RedisHtmlCache) GetHtml(key string) (string, error) {
	cmd := x.redisClient.HGet(context.Background(), "xxx:HTML", key)
	r, err := cmd.Result()
	if err != nil {
		if err == rV8.Nil {
			return "", nil // key不存在，不当作错误
		}
		return "", serr.WithStack(err)
	}
	return r, nil
}

func (x *RedisHtmlCache) Exists(key string) (bool, error) {
	cmd := x.redisClient.HExists(context.Background(), "xxx:HTML", key)
	r, err := cmd.Result()
	if err != nil {
		return false, serr.WithStack(err)
	}
	return r, nil
}

func (x *RedisHtmlCache) SetHtml(key, value string) error {
	err := x.redisClient.HSet(context.Background(), "xxx:HTML", key, value).Err()
	if err != nil {
		return serr.WithStack(err)
	}
	return nil
}
