package cms

import (
	"regexp"
	"strings"
	"sync"

	"github.com/CloudyKit/jet/v6"
	"github.com/syncfuture/go/shttp"
	"github.com/syncfuture/go/spool"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

type ICMS interface {
	GetContent(key string, args ...interface{}) string
	Render(key string, args ...interface{}) (string, error)
}

const _err1 = "could not be found"

var (
	_paramPool = sync.Pool{
		New: func() interface{} {
			return make(jet.VarMap)
		},
	}
	_bufferPool = spool.NewSyncBufferPool(1024)
	_minifier   = minify.New()
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

func GetParams() jet.VarMap {
	return _paramPool.Get().(jet.VarMap)
}

func ReleaseParams(params jet.VarMap) {
	for k := range params {
		delete(params, k)
	}
	_paramPool.Put(params)
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
