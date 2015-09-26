package factory

import (
	"github.com/blevesearch/bleve"
	"os"
)

type BleveEngine struct {
}

/* Implements the SearchEngine interface */

func (be *BleveEngine) Index(document []byte, id string) error {
	var index bleve.Index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(INDEX, mapping)
	if err != nil {
		index, _ = bleve.Open(INDEX)
	}
	index.Index(id, document)
	index.Close()
	return nil
}

func (be *BleveEngine) Search(query string) (interface{}, error) {
	index, _ := bleve.Open(INDEX)
	bleveQuery := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(bleveQuery)
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	return searchResults, nil
}

func (be *BleveEngine) Delete() error {
	index, err := bleve.Open(INDEX)
	if err != nil {
		return err
	}
	index.NewBatch().Reset()
	index.Close()
	os.RemoveAll(INDEX)
	return nil
}
