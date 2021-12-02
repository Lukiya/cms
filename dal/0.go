package dal

type IHtmlCacheDAL interface {
	GetHtml(path string) (string, error)
	SetHtml(path, value string) error
}
