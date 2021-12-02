package dal

import "github.com/syncfuture/go/u"

var (
	_dal IHtmlDAL
)

func init() {
	_dal = newRedisHtmlDAL()
}

func GetHtml(key string) string {
	r, err := _dal.GetHtml(key)
	u.LogError(err)
	return r
}
