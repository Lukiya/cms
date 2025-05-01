package redis

import (
	"context"

	"github.com/DreamvatLab/go/xerr"
	"github.com/DreamvatLab/go/xredis"
	"github.com/Lukiya/cms/dal"
	rV9 "github.com/redis/go-redis/v9"
)

func NewRedisContentDAL(redisKey string, redisConfig *xredis.RedisConfig) dal.IContentDAL {
	r := new(RedisHtmlCache)
	r.redisClient = xredis.NewClient(redisConfig)
	r.redisKey = redisKey
	return r
}

type RedisHtmlCache struct {
	redisBase
}

func (x *RedisHtmlCache) GetContent(path string) (string, error) {
	cmd := x.redisClient.HGet(context.Background(), x.redisKey, path)
	r, err := cmd.Result()
	if err != nil {
		if err == rV9.Nil {
			return "", nil // key不存在，不当作错误
		}
		return "", xerr.WithStack(err)
	}
	return r, nil
}

func (x *RedisHtmlCache) Exists(path string) (bool, error) {
	cmd := x.redisClient.HExists(context.Background(), x.redisKey, path)
	r, err := cmd.Result()
	if err != nil {
		return false, xerr.WithStack(err)
	}
	return r, nil
}

func (x *RedisHtmlCache) SetContent(path, value string) error {
	err := x.redisClient.HSet(context.Background(), x.redisKey, path, value).Err()
	if err != nil {
		return xerr.WithStack(err)
	}
	return nil
}
