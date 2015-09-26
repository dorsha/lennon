package factory

const (
	INDEX          = "index"
	VENDOR_ELASTIC = "elastic"
	VENDOR_BLEVE   = "bleve"
)

type Document struct {
	Id   string
	Data []byte
}

type SearchEngine interface {
	BatchIndex(documents *[]*Document) (int64, error)
	Index(document *Document) (int64, error)
	Search(query string) (interface{}, error)
	Delete() error
}
