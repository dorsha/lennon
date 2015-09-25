package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

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
	url       = flag.String("url", "", "Search engine URL (i.e. http://192.168.1.26:9200)")
)

func main() {
	flag.Parse()

	fmt.Printf("Search vendor: %s\n", *vendor)

	if *action == "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", "Action must be specified.")
		os.Exit(1)
	}

	if *url == "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", "Search engine URL must be specified.")
		os.Exit(1)
	}

	switch *action {
	case index:
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

	// Print result
	if elasticSearchResult, ok := searchResult.(*elastic.SearchResult); ok {

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
}

func getSearchEngine() factory.SearchEngine {
	engine, err := factory.GetSearchEngine(url, vendor)
	utils.ErrorCheck(err)
	return engine
}
