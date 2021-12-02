package main

import (
	"github.com/Lukiya/cms/core"
	"github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/valyala/fasthttp"
)

func main() {
	// core.Host.GET("/", home.Index)

	// core.Host.ServeFiles("/{filepath:*}", "www")
	fasthttp.ListenAndServe(":8000", fastHTTPHandler)
	slog.Fatal(core.Host.Run())
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	a := ctx.URI().RequestURI()
	print(u.BytesToStr(a))
}
