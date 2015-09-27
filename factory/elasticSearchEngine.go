package factory

import (
	"log"
	"os"
	"time"

	"github.com/dorsha/lennon/utils"
	"gopkg.in/olivere/elastic.v2"
)

type ElasticEngine struct {
	Client *elastic.Client
}

/* Implements the SearchEngine interface */

func (es *ElasticEngine) BatchIndex(documents []*Document) (int64, error) {
	err := es.createIndexIfNotExists()
	utils.ErrorCheck(err)

	bulkRequest := es.Client.Bulk()

	for _, document := range documents {
		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(INDEX).Type(document.Id).Id(document.Id).
			Doc(string((*document).Data[:])))
	}

	bulkResponse, err := bulkRequest.Do()

	return int64(bulkResponse.Took), err
}

func (es *ElasticEngine) Index(document *Document) (int64, error) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	err := es.createIndexIfNotExists()
	utils.ErrorCheck(err)
	// Index the data
	_, err = es.Client.Index().Index(INDEX).Type((*document).Id).Id((*document).Id).
		BodyJson(string((*document).Data[:])).Do()
	if err != nil {
		return 0, err
	}

	return int64(time.Now().UnixNano()/int64(time.Millisecond) - start), nil
}

func (es *ElasticEngine) Search(query string) (interface{}, error) {
	termQuery := elastic.NewQueryStringQuery(query)
	searchResult, err := es.Client.Search().
		Index(INDEX).
		Query(&termQuery).
		From(0).Size(10).
		Pretty(true).
		Do()

	return searchResult, err
}

func (es *ElasticEngine) Delete() error {
	_, err := es.Client.DeleteIndex(INDEX).Do()
	return err
}

func CreateElasticClient(url *string) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(*url),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
}

func (es *ElasticEngine) createIndexIfNotExists() error {
	exists, _ := es.Client.IndexExists(INDEX).Do()

	if !exists {
		_, err := es.Client.CreateIndex(INDEX).Do()
		if err != nil {
			return err
		}
	}

	return nil
}
