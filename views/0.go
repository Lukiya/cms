package views

import (
	"bytes"

	"github.com/syncfuture/go/shttp"
	"github.com/syncfuture/go/spool"
	"github.com/syncfuture/host"
)

var (
	_bfPool    = spool.NewSyncBufferPool(1024)
	_modelPool = spool.NewSyncMapPool()
)

type IPageRender func(model *map[string]interface{}, buffer *bytes.Buffer)

func GetModel() *map[string]interface{} {
	return _modelPool.GetMap()
}

// func PutModel(m *map[string]interface{}) {
// 	_modelPool.PutMap(m)
// }

func Render(ctx host.IHttpContext, pageRender IPageRender, models ...*map[string]interface{}) {
	var model *map[string]interface{}
	if len(models) == 1 {
		model = models[0]
	} else {
		model = GetModel()
	}

	// (*model)["Debug"] = core.Host.GetDebug()
	// (*model)["IsAuthenticated"] = host.GetUserID(ctx, core.Host.GetUserIDSessionKey()) != ""
	(*model)["Debug"] = true
	(*model)["IsAuthenticated"] = true

	bf := _bfPool.GetBuffer()
	defer func() {
		_bfPool.PutBuffer(bf)
		_modelPool.PutMap(model)
	}()

	pageRender(model, bf)

	ctx.SetContentType(shttp.CTYPE_HTML)
	bf.WriteTo(ctx)
}

func boolText(flag bool, strTrue string, strFalse string) string {
	if flag {
		return strTrue
	} else {
		return strFalse
	}
}
