package dal

type IHtmlCacheDAL interface {
	GetHtml(key string) (string, error)
	SetHtml(key, value string) error
}
