package cms

import (
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/Lukiya/cms/dal"
	"github.com/Lukiya/cms/dal/redis"
	"github.com/syncfuture/go/sconfig"
	"github.com/syncfuture/go/serr"
	"github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/u"
)

type IJetCMS interface {
	ICMS
	GetViewEngine() *jet.Set
}

func NewJetCMS(cp sconfig.IConfigProvider) IJetCMS {
	var redisConfig *sredis.RedisConfig
	cp.GetStruct("Redis", &redisConfig)
	htmlCacheKey := cp.GetString("CMS.ContentKey")
	templateKey := cp.GetString("CMS.TemplateKey")
	isDebug := cp.GetBool("Debug")

	var viewEngine *jet.Set
	if isDebug {
		// 调试
		viewEngine = jet.NewSet(
			redis.NewRedisTemplateLoader(templateKey, redisConfig),
			jet.InDevelopmentMode(),
		)
	} else {
		// 生产
		viewEngine = jet.NewSet(
			redis.NewRedisTemplateLoader(templateKey, redisConfig),
		)
	}

	return &jetCMS{
		cp:         cp,
		htmlCache:  redis.NewRedisContentDAL(htmlCacheKey, redisConfig),
		viewEngine: viewEngine,
	}
}

type jetCMS struct {
	htmlCache  dal.IContentDAL
	viewEngine *jet.Set
	cp         sconfig.IConfigProvider
}

func (x *jetCMS) GetViewEngine() *jet.Set {
	return x.viewEngine
}
func (x *jetCMS) GetContent(key string, args ...interface{}) string {
	r, err := x.htmlCache.GetContent(key) // 查找缓存
	if u.LogError(err) {
		return ""
	}

	if r == "" { // 缓存为空
		r, err = x.Render(key, args...) // 渲染
		if err != nil {
			if strings.Contains(err.Error(), _err1) { // 没有找到模板，只记录debug信息
				slog.Debug(err.Error())
			} else {
				slog.Error(err)
			}
		}

		if r != "" {
			if len(args) > 1 {
				if ok := args[1].(bool); ok {
					ctype := GetContentType(key)
					r1, err := _minifier.String(ctype, r)
					if !u.LogError(err) {
						r = r1
					}
				}
			}

			err = x.htmlCache.SetContent(key, r) // 放入缓存
			u.LogError(err)
		}
	}

	return r
}

func (x *jetCMS) Render(key string, args ...interface{}) (string, error) {
	var params jet.VarMap

	if len(args) > 0 && args[0] != nil { // 第一个参数作为 jet 数据模型
		defer releaseParams(params) // 使用完毕释放
		params = args[0].(jet.VarMap)
	}

	template, err := x.viewEngine.GetTemplate(key) // 获取模板
	if err != nil {
		return "", serr.WithStack(err)
	}

	r := _bufferPool.GetBuffer()
	defer _bufferPool.PutBuffer(r) // buffer使用完毕释放

	err = template.Execute(r, params, nil) // 执行渲染
	if err != nil {
		return "", serr.WithStack(err)
	}
	return u.BytesToStr(r.Bytes()), nil
}
