package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/dorsha/lennon/factory"
	"github.com/dorsha/lennon/utils"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"os"
)

const (
	index  = "index"
	search = "search"
)

var (
	vendor = flag.String("vendor", "elastic", "Specify the vendor you want to test ("+factory.VENDOR_ELASTIC+"/"+
		factory.VENDOR_BLEVE+")")
	action    = flag.String("action", "", "What do you want to do? Avialble actions: "+index+","+search)
	pathToDoc = flag.String("document", "", "Path to document that you want to index")
	query     = flag.String("query", "", "Search query")
	url       = flag.String("url", "", "ElasticSearch URL (i.e. http://192.168.1.26:9200) - not relvant for Bleve")
)

func main() {
	flag.Parse()

	fmt.Printf("Search vendor: %s\n", *vendor)

	if *action == "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", "Action must be specified.")
		os.Exit(1)
	}

	switch *action {
	case index:
		if *vendor == factory.VENDOR_ELASTIC && *url == "" {
			fmt.Fprintf(os.Stderr, "error: %v\n", "Search engine URL must be specified when using ElasticSearch.")
			os.Exit(1)
		}
		doIndex()
	case search:
		doSearch()
	default:
		fmt.Fprintf(os.Stderr, "error: %v\n", "Action must be one of %s or %s.", index, search)
		os.Exit(1)
	}
}

func doIndex() {
	engine := getSearchEngine()

	if *pathToDoc == "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", "Path to document must be specified")
		os.Exit(1)
	}

	doc, err := ioutil.ReadFile(*pathToDoc)
	utils.ErrorCheck(err)

	start := time.Now().UnixNano() / int64(time.Millisecond)
	err = engine.Index(doc)
	utils.ErrorCheck(err)
	fmt.Printf("Took %d milliseconds to index.\n", time.Now().UnixNano()/int64(time.Millisecond)-start)
}

func doSearch() {

	if *query == "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", "Query string must be specified")
		os.Exit(1)
	}

	engine := getSearchEngine()

	// Search
	searchResult, err := engine.Search(*query)
	utils.ErrorCheck(err)

	switch value := searchResult.(type) {
	case *elastic.SearchResult:
		printElasticResult(value)
	case *bleve.SearchResult:
		printBleveResult(value)
	}
}

func printElasticResult(elasticSearchResult *elastic.SearchResult) {
	fmt.Printf("Query took %d milliseconds\n", elasticSearchResult.TookInMillis)
	fmt.Printf("Query hits: %d\n", elasticSearchResult.TotalHits())

	fmt.Printf("Query Result:\n")
	if elasticSearchResult.Hits != nil {
		for _, hit := range elasticSearchResult.Hits.Hits {
			var p interface{}
			json.Unmarshal(*hit.Source, &p)
			fmt.Println(p)
		}
	}
}

func printBleveResult(bleveSearchResult *bleve.SearchResult) {
	fmt.Printf("Query took %d milliseconds\n", bleveSearchResult.Took)
	fmt.Printf("Query hits: %d\n", bleveSearchResult.Total)
	fmt.Printf("Query Result:\n")
	fmt.Println(bleveSearchResult)
}

func getSearchEngine() factory.SearchEngine {
	engine, err := factory.GetSearchEngine(url, vendor)
	utils.ErrorCheck(err)
	return engine
}
