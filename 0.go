package cms

import (
	"strings"
	"sync"

	"github.com/CloudyKit/jet/v6"
	"github.com/syncfuture/go/shttp"
	"github.com/syncfuture/go/spool"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
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
	_minifier.AddFunc(shttp.CTYPE_HTML, html.Minify)
	_minifier.AddFunc(shttp.CTYPE_JS, js.Minify)

	// _minifier.AddFunc("image/svg+xml", svg.Minify)
}

func MakeParams() jet.VarMap {
	return _paramPool.Get().(jet.VarMap)
}

func releaseParams(params jet.VarMap) {
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
	}
	return shttp.CTYPE_HTML
}
