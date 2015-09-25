package factory

import (
	"errors"
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
}

type Engine struct {
	client *elastic.Client
}

type ElasticEngine struct {
	Engine
}

type BleveEngine struct {
	Engine
}

type es ElasticEngine
type bleve BleveEngine

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

func (es *BleveEngine) Index(document []byte) error {
	// todo - implement bleve
	return nil
}

func (es *BleveEngine) Search(query string) (interface{}, error) {
	// todo - implement bleve
	return nil, nil
}

func GetSearchEngine(url *string, vendor *string) (SearchEngine, error) {
	// Create a client
	client, err := createElasticClient(url)
	if err != nil {
		return nil, err
	}

	var engine SearchEngine
	switch *vendor {
	case VENDOR_ELASTIC:
		engine = &ElasticEngine{Engine{client}}
	case VENDOR_BLEVE:
		engine = &BleveEngine{Engine{client}}
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
