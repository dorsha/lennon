package factory

import (
	"errors"
	"github.com/blevesearch/bleve"
	"gopkg.in/olivere/elastic.v2"
	"log"
	"os"
	"time"
)

const (
	INDEX          = "doc.test"
	VENDOR_ELASTIC = "elastic"
	VENDOR_BLEVE   = "bleve"
)

type SearchEngine interface {
	Index(document []byte) error
	Search(query string) (interface{}, error)
	Delete() error
}

type ElasticEngine struct {
	client *elastic.Client
}

type BleveEngine struct {
}

func (es *ElasticEngine) Index(document []byte) error {
	// create index if not exists
	exists, err := es.client.IndexExists(INDEX).Do()
	if !exists {
		_, err := es.client.CreateIndex(INDEX).Do()
		if err != nil {
			return err
		}
	}

	// Index the data
	_, err = es.client.Index().Index(INDEX).Type("person").Id("1").BodyJson(string(document[:])).Do()
	if err != nil {
		return err
	}

	return nil
}

func (es *ElasticEngine) Search(query string) (interface{}, error) {
	termQuery := elastic.NewQueryStringQuery(query)
	searchResult, err := es.client.Search().
		Index(INDEX).
		Query(&termQuery).
		From(0).Size(10).
		Pretty(true).
		Sort("age", true).
		Do()

	return searchResult, err
}

func (es *ElasticEngine) Delete() error {
	_, err := es.client.DeleteIndex(INDEX).Do()
	return err
}

func (be *BleveEngine) Index(document []byte) error {
	var index bleve.Index
	index, err := bleve.Open(INDEX)
	if index == nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(INDEX, mapping)
		if err != nil {
			return err
		}
	}

	index.Index("data", document)
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
	index, _ := bleve.Open(INDEX)
	return index.Delete(INDEX)
}

func GetSearchEngine(url *string, vendor *string) (SearchEngine, error) {

	var engine SearchEngine
	switch *vendor {
	case VENDOR_ELASTIC:
		// Create a client
		client, err := createElasticClient(url)
		if err != nil {
			return nil, err
		}
		engine = &ElasticEngine{client}
	case VENDOR_BLEVE:
		engine = &BleveEngine{}
	default:
		return nil, errors.New("Engine vendor must be specified.")

	}

	return engine, nil
}

func createElasticClient(url *string) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(*url),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
}
