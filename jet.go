package cms

import (
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/Lukiya/cms/dal"
	"github.com/Lukiya/cms/dal/redis"
	"github.com/syncfuture/go/sconfig"
	"github.com/syncfuture/go/serr"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/u"
)

type IJetHtmlCMS interface {
	ICMS
	GetViewEngine() *jet.Set
}

func NewJetHtmlCMS(cp sconfig.IConfigProvider) IJetHtmlCMS {
	var redisConfig *sredis.RedisConfig
	cp.GetStruct("Redis", &redisConfig)
	htmlCacheKey := cp.GetString("CMS.HtmlCacheKey")
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

	return &jetHtmlCMS{
		cp:         cp,
		htmlCache:  redis.NewRedisHtmlCache(htmlCacheKey, redisConfig),
		viewEngine: viewEngine,
	}
}

type jetHtmlCMS struct {
	htmlCache  dal.IHtmlCacheDAL
	viewEngine *jet.Set
	cp         sconfig.IConfigProvider
}

func (x *jetHtmlCMS) GetViewEngine() *jet.Set {
	return x.viewEngine
}
func (x *jetHtmlCMS) GetHtml(key string, args ...interface{}) string {
	r, err := x.htmlCache.GetHtml(key) // 查找缓存
	if u.LogError(err) {
		return ""
	}

	if r == "" { // 缓存为空
		r, err = x.Render(key, args...) // 渲染
		if u.LogError(err) {
			return ""
		}

		err = x.htmlCache.SetHtml(key, r) // 放入缓存
		u.LogError(err)
	}

	return r
}

func (x *jetHtmlCMS) Render(key string, args ...interface{}) (string, error) {
	var params jet.VarMap

	if len(args) > 0 { // 第一个参数作为 jet 数据模型
		defer releaseParams(params) // 使用完毕释放
		params = args[0].(jet.VarMap)
	}

	template, err := x.viewEngine.GetTemplate(key) // 获取模板
	if err != nil {
		if strings.Contains(err.Error(), "not found") { // 没有找到模板，只记录debug信息
			log.Debug(err.Error())
			return "", nil
		}
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
