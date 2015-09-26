package factory

const (
	INDEX          = "index"
	VENDOR_ELASTIC = "elastic"
	VENDOR_BLEVE   = "bleve"
)

type SearchEngine interface {
	Index(document []byte, id string) error
	Search(query string) (interface{}, error)
	Delete() error
}
