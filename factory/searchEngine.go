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
	BatchIndex(documents *[]*Document) error
	Index(document *Document) error
	Search(query string) (interface{}, error)
	Delete() error
}
