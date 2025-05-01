package jet

import (
	"errors"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/DreamvatLab/go/xbytes"
	"github.com/DreamvatLab/go/xconfig"
	"github.com/DreamvatLab/go/xerr"
	"github.com/DreamvatLab/go/xlog"
	"github.com/DreamvatLab/go/xredis"
	"github.com/DreamvatLab/go/xsync"
	"github.com/Lukiya/cms"
	"github.com/Lukiya/cms/dal"
	"github.com/Lukiya/cms/dal/redis"
	"github.com/tdewolff/minify/v2"
)

const _err1 = "could not be found"

var (
	_bufferPool = xsync.NewSyncBufferPool(1024)
	_minifier   = minify.New()
)

type IJetCMS interface {
	cms.ICMS
	GetViewEngine() *jet.Set
}

func NewJetCMS(cp xconfig.IConfigProvider) IJetCMS {
	var redisConfig *xredis.RedisConfig
	cp.GetStruct("Redis", &redisConfig)

	htmlCacheKey := cp.GetString("CMS.ContentKey")
	templateKey := cp.GetString("CMS.TemplateKey")

	// Loader
	loaderStoreProvider := cp.GetString("CMS.LoaderStore.Provider")
	if loaderStoreProvider != "File" && loaderStoreProvider != "Redis" {
		xerr.FatalIfErr(errors.New("CMS.LoaderStore.Provider can only be ether 'File' or 'Redis'"))
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

func NewCMS(cp xconfig.IConfigProvider) cms.ICMS {
	return NewJetCMS(cp)
}

type jetCMS struct {
	htmlCache  dal.IContentDAL
	viewEngine *jet.Set
	cp         xconfig.IConfigProvider
}

func (x *jetCMS) GetViewEngine() *jet.Set {
	return x.viewEngine
}

// GetContent 获取内容
// key: 内容键，一般为路径
// arg[0]: *jet.GetParams获取, 使用完毕用jet.ReleaseParams释放
// arg[1]: bool 是否使用缓存
// arg[2]: bool 输出html是否压缩
func (x *jetCMS) GetContent(key string, args ...interface{}) string {
	cache := false
	if len(args) > 1 { // 判断是否使用缓存
		cache = args[1].(bool)
	}

	if cache {
		// 使用缓存
		r, err := x.htmlCache.GetContent(key) // 查找缓存
		if xerr.LogError(err) {
			return ""
		}

		if r == "" { // 缓存为空
			r = x.render(key, args...) // 渲染
			if r != "" {
				err = x.htmlCache.SetContent(key, r) // 放入缓存
				xerr.LogError(err)
			}
		}

		return r
	}

	// 直接渲染
	return x.render(key, args...)
}

// render 渲染内容
// key: 内容键，一般为路径
// arg[0]: *JetParams，用jet.GetParams获取, 使用完毕用jet.ReleaseParams释放
// arg[1]: bool 是否使用缓存
// arg[2]: bool 输出html是否minify
func (x *jetCMS) render(key string, args ...interface{}) string {
	r, err := x.Render(key, args...) // 渲染
	if err != nil {
		if strings.Contains(err.Error(), _err1) { // 没有找到模板，只记录debug信息，不当成错误
			xlog.Debug(err.Error())
		} else {
			xlog.Error(err)
		}
	}

	if r != "" && len(args) > 2 {
		// 是否 minify
		if ok := args[2].(bool); ok {
			ctype := cms.GetContentType(key)
			r1, err := _minifier.String(ctype, r)
			if !xerr.LogError(err) {
				r = r1
			}
		}
	}

	return r
}

// Render 渲染内容
// key: 内容键，一般为路径
// arg[0]: *JetParams，用jet.GetParams获取, 使用完毕用jet.ReleaseParams释放
// arg[1]: bool 是否使用缓存
// arg[2]: bool 输出html是否压缩
func (x *jetCMS) Render(key string, args ...interface{}) (string, error) {
	var params *JetParams

	if len(args) > 0 && args[0] != nil { // 第一个参数作为 jet 数据模型
		params = args[0].(*JetParams)
	}

	template, err := x.viewEngine.GetTemplate(key) // 获取模板
	if err != nil {
		return "", xerr.WithStack(err)
	}

	r := _bufferPool.GetBuffer()
	defer _bufferPool.PutBuffer(r) // buffer使用完毕释放

	err = template.Execute(r, *params.data, nil) // 执行渲染
	if err != nil {
		return "", xerr.WithStack(err)
	}
	return xbytes.BytesToStr(r.Bytes()), nil
}
