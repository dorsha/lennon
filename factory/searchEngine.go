package factory

const (
	INDEX          = "id"
	VENDOR_ELASTIC = "elastic"
	VENDOR_BLEVE   = "bleve"
)

type SearchEngine interface {
	Index(document []byte) error
	Search(query string) (interface{}, error)
	Delete() error
}
