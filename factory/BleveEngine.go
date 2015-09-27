package factory

import (
	"fmt"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/store/boltdb"
	"github.com/blevesearch/bleve/index/store/goleveldb"
)

type BleveEngine struct {
}

/* Implements the SearchEngine interface */

func (be *BleveEngine) BatchIndex(documents []*Document) (int64, error) {
	start := time.Now().UnixNano() / int64(time.Millisecond)
	var index bleve.Index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(INDEX, mapping)
	if err != nil {
		index, _ = bleve.Open(INDEX)
	}

	batch := index.NewBatch()

	for _, document := range documents {
		batch.Index(document.Id, document.Data)
	}

	index.Batch(batch)
	index.Close()

	return time.Now().UnixNano()/int64(time.Millisecond) - start, nil
}

func (be *BleveEngine) Index(document *Document) (int64, error) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	var index bleve.Index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(INDEX, mapping)
	if err != nil {
		index, _ = bleve.Open(INDEX)
	}
	id := (*document).Id
	data := (*document).Data
	index.Index(id, data)
	index.Close()

	return time.Now().UnixNano()/int64(time.Millisecond) - start, nil
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
	index.Close()
	os.RemoveAll(INDEX)
	return nil
}

func (be *BleveEngine) SetKVStore(storeName string) error {
	if storeName != goleveldb.Name && storeName != boltdb.Name {
		panic("Unsupprted KV store.")
	}

	bleve.Config.DefaultKVStore = storeName
	fmt.Printf("KV Store: %s\n", storeName)
	return nil
}
