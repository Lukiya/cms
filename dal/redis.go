package dal

import (
	"github.com/Lukiya/cms/core"
	"github.com/go-redis/redis/v8"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/u"
)

var (
	redisClient redis.UniversalClient
)

func init() {
	var redisConfig *sredis.RedisConfig
	err := core.Config.GetStruct("Redis", &redisConfig)
	u.LogFaltal(err)
	redisClient = sredis.NewClient(redisConfig)
}

func newRedisHtmlDAL() IHtmlDAL {
	return new(RedisHtmlDAL)
}

type RedisHtmlDAL struct{}

func (x *RedisHtmlDAL) GetHtml(key string) (string, error) {
	return "1111", nil
}
