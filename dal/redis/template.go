package redis

import (
	"bytes"
	"context"
	"io"

	"github.com/CloudyKit/jet/v6"
	"github.com/DreamvatLab/go/xerr"
	"github.com/DreamvatLab/go/xredis"
)

type RedisTemplateLoader struct {
	redisBase
}

func NewRedisTemplateLoader(key string, redisConfig *xredis.RedisConfig) jet.Loader {
	r := new(RedisTemplateLoader)
	r.redisClient = xredis.NewClient(redisConfig)
	r.redisKey = key
	return r
}

func (x *RedisTemplateLoader) Exists(templatePath string) bool {
	ex := x.redisClient.HExists(context.Background(), x.redisKey, templatePath)
	r, err := ex.Result()
	if xerr.LogError(err) {
		return false
	}

	return r
}

func (x *RedisTemplateLoader) Open(templatePath string) (io.ReadCloser, error) {
	s := x.redisClient.HGet(context.Background(), x.redisKey, templatePath)
	b, err := s.Bytes()
	if err != nil {
		return nil, xerr.WithStack(err)
	}
	return io.NopCloser(bytes.NewReader(b)), nil
}
