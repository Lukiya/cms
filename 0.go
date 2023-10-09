package cms

import (
	"regexp"
	"strings"

	"github.com/syncfuture/go/shttp"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

type ICMS interface {
	// GetContent 获取内容
	// key: 内容键，一般为路径
	// arg[0]: 渲染使用的页面参数数据
	// arg[1]: bool 是否使用缓存
	// arg[2]: bool 输出html是否压缩
	GetContent(key string, args ...interface{}) string
	// GetContent 获取内容
	// key: 内容键，一般为路径
	// arg[0]: 渲染使用的页面参数数据
	// arg[1]: bool 是否使用缓存
	// arg[2]: bool 输出html是否压缩
	Render(key string, args ...interface{}) (string, error)
}

var (
	// _paramPool = sync.Pool{
	// 	New: func() interface{} {
	// 		return make(jet.VarMap)
	// 	},
	// }
	// _bufferPool = spool.NewSyncBufferPool(1024)
	_minifier = minify.New()
)

func init() {
	_minifier.AddFunc(shttp.CTYPE_CSS, css.Minify)
	_minifier.Add(shttp.CTYPE_HTML, &html.Minifier{
		KeepComments:            false,
		KeepConditionalComments: false,
		KeepDefaultAttrVals:     false,
		KeepDocumentTags:        false,
		KeepEndTags:             false,
		KeepQuotes:              true,
		KeepWhitespace:          false,
	})
	_minifier.AddFunc("image/svg+xml", svg.Minify)
	_minifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	_minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	_minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

}

func GetContentType(path string) string {
	if strings.HasSuffix(path, ".css") {
		return shttp.CTYPE_CSS
	} else if strings.HasSuffix(path, ".js") {
		return shttp.CTYPE_JS
	} else if strings.HasSuffix(path, ".xml") {
		return shttp.CTYPE_XML
	} else if strings.HasSuffix(path, ".json") {
		return shttp.CTYPE_JSON
	}
	return shttp.CTYPE_HTML
}
