package factory

import (
	"errors"
	"gopkg.in/olivere/elastic.v2"
	"log"
	"os"
	"time"
)

type ElasticEngine struct {
	client *elastic.Client
}

/* Implements the SearchEngine interface */

func (es *ElasticEngine) Index(document []byte, id string) error {
	// create index if not exists
	exists, err := es.client.IndexExists(INDEX).Do()

	if !exists {
		_, err := es.client.CreateIndex(INDEX).Do()
		if err != nil {
			return err
		}
	}

	// Index the data
	_, err = es.client.Index().Index(INDEX).Type(id).Id(id).BodyJson(string(document[:])).Do()
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
		Do()

	return searchResult, err
}

func (es *ElasticEngine) Delete() error {
	_, err := es.client.DeleteIndex(INDEX).Do()
	return err
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
