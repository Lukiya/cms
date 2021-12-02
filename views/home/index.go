package home

import (
	"github.com/Lukiya/cms/views"
	"github.com/syncfuture/host"
)

func Index(ctx host.IHttpContext) {
	model := views.GetModel()
	views.Render(ctx, views.HomeIndex, model)

}
