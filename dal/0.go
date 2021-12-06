package dal

type IContentDAL interface {
	GetContent(path string) (string, error)
	SetContent(path, value string) error
}
