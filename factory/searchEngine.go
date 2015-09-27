package factory

import "errors"

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
	BatchIndex(documents []*Document) (int64, error)
	Index(document *Document) (int64, error)
	Search(query string) (interface{}, error)
	Delete() error
}

func GetSearchEngine(url *string, vendor *string, KVStore string) (SearchEngine, error) {
	var engine SearchEngine
	switch *vendor {
	case VENDOR_ELASTIC:
		// Create a client
		client, err := CreateElasticClient(url)
		if err != nil {
			return nil, err
		}
		engine = &ElasticEngine{client}
	case VENDOR_BLEVE:
		bleveEngine := &BleveEngine{}
		bleveEngine.SetKVStore(KVStore)
		engine = bleveEngine
	default:
		return nil, errors.New("Engine vendor must be specified.")
	}

	return engine, nil
}
