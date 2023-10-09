package cms

import (
	"errors"
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

	// Loader
	loaderStoreProvider := cp.GetString("CMS.LoaderStore.Provider")
	if loaderStoreProvider != "File" && loaderStoreProvider != "Redis" {
		u.LogFatal(errors.New("CMS.LoaderStore.Provider can only be ether 'File' or 'Redis'"))
	}
	viewsFileDirPath := cp.GetString("CMS.LoaderStore.File.DirPath")
	if viewsFileDirPath == "" {
		viewsFileDirPath = "./views"
	}
	var loaderStore jet.Loader
	if loaderStoreProvider == "File" {
		loaderStore = jet.NewOSFileSystemLoader(viewsFileDirPath)
	} else if loaderStoreProvider == "Redis" {
		loaderStore = redis.NewRedisTemplateLoader(templateKey, redisConfig)
	}

	isDebug := cp.GetBool("Debug")

	var viewEngine *jet.Set
	if isDebug {
		// 调试
		viewEngine = jet.NewSet(
			loaderStore,
			jet.InDevelopmentMode(),
		)
	} else {
		// 生产
		viewEngine = jet.NewSet(
			loaderStore,
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

// GetContent 获取内容，args[0]指定是否使用缓存
func (x *jetCMS) GetContent(key string, args ...interface{}) string {
	cache := false
	if len(args) > 1 { // 判断是否使用缓存
		cache = args[1].(bool)
	}

	if cache {
		// 使用缓存
		r, err := x.htmlCache.GetContent(key) // 查找缓存
		if u.LogError(err) {
			return ""
		}

		if r == "" { // 缓存为空
			r = x.render(key, args...) // 渲染
			if r != "" {
				err = x.htmlCache.SetContent(key, r) // 放入缓存
				u.LogError(err)
			}
		}

		return r
	}

	// 直接渲染
	return x.render(key, args...)
}

func (x *jetCMS) render(key string, args ...interface{}) string {
	r, err := x.Render(key, args...) // 渲染
	if err != nil {
		if strings.Contains(err.Error(), _err1) { // 没有找到模板，只记录debug信息，不当成错误
			slog.Debug(err.Error())
		} else {
			slog.Error(err)
		}
	}

	if r != "" && len(args) > 2 {
		// 是否 minify
		if ok := args[2].(bool); ok {
			ctype := GetContentType(key)
			r1, err := _minifier.String(ctype, r)
			if !u.LogError(err) {
				r = r1
			}
		}
	}

	return r
}

func (x *jetCMS) Render(key string, args ...interface{}) (string, error) {
	var params jet.VarMap

	if len(args) > 0 && args[0] != nil { // 第一个参数作为 jet 数据模型
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
