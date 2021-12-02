package dal

type IHtmlDAL interface {
	GetHtml(key string) (string, error)
}
