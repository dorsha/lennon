package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"io/ioutil"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/store/boltdb"
	"github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/dorsha/lennon/factory"
	"github.com/dorsha/lennon/utils"
	"gopkg.in/olivere/elastic.v2"
)

const (
	index       = "index"
	search      = "search"
	deleteIndex = "deleteIndex"
)

var (
	vendor = flag.String("vendor", "elastic", "Specify the vendor you want to test ("+factory.VENDOR_ELASTIC+"/"+
		factory.VENDOR_BLEVE+")")
	action       = flag.String("action", "", "What do you want to do? Avialble actions: "+index+", "+search+", "+deleteIndex)
	pathToDoc    = flag.String("document", "", "Path to document that you want to index")
	pathToFolder = flag.String("folder", "", "Path to folder that contains documents to index (non-recursive)")
	query        = flag.String("query", "", "Search query")
	url          = flag.String("url", "", "ElasticSearch URL (i.e. http://192.168.1.26:9200) - not relvant for Bleve")
	KVStore      = flag.String("store", goleveldb.Name, "KV Store for Bleve. Avialble stores: "+goleveldb.Name+
		", "+boltdb.Name)
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
	case deleteIndex:
		doDelete()
	default:
		fmt.Fprintf(os.Stderr, "error: %v\n", "Invalid action")
		os.Exit(1)
	}
}

func doIndex() {
	engine := getSearchEngine()

	if *pathToDoc == "" && *pathToFolder == "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", "Path to document/folder must be specified")
		os.Exit(1)
	}

	if *pathToDoc != "" && *pathToFolder != "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", "Folder or Document? Can't specify both.")
		os.Exit(1)
	}

	if *pathToDoc != "" {
		timeTook := indexDocument(engine)
		fmt.Printf("Took %d milliseconds to index.\n", timeTook)
	} else {
		timeTook, docCount := indexDocuments(engine)
		fmt.Printf("Took %d milliseconds to index %d documents.\n", timeTook, docCount)
	}
}

func indexDocument(engine factory.SearchEngine) (timeTook int64) {
	data, err := utils.ReadFile(*pathToDoc)
	utils.ErrorCheck(err)
	fmt.Printf("Indexing document: %s\n", *pathToDoc)
	doc := factory.Document{utils.FixIdSyntax(*pathToDoc), data}
	timeTook, err = engine.Index(&doc)
	utils.ErrorCheck(err)
	return timeTook
}

func indexDocuments(engine factory.SearchEngine) (int64, int) {
	// prepare documents for batch index
	files, err := ioutil.ReadDir(*pathToFolder)
	utils.ErrorCheck(err)
	var docCount int = 0
	var documents = make([]*factory.Document, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue // expected flat structure
		}
		path := *pathToFolder + "/" + file.Name()
		data, err := utils.ReadFile(path)
		utils.ErrorCheck(err)
		documents[docCount] = &factory.Document{utils.FixIdSyntax(path), data}
		docCount++
	}
	// index in batch
	timeTook, err := engine.BatchIndex(documents)
	utils.ErrorCheck(err)
	return timeTook, docCount
}

func doDelete() {
	engine := getSearchEngine()
	err := engine.Delete()
	utils.ErrorCheck(err)
	fmt.Printf("Index deleted\n")
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
	fmt.Printf("Query took %d nanoseconds\n", bleveSearchResult.Took)
	fmt.Printf("Query hits: %d\n", bleveSearchResult.Total)
	fmt.Printf("Query Result:\n")
	fmt.Println(bleveSearchResult)
}

func getSearchEngine() factory.SearchEngine {
	engine, err := factory.GetSearchEngine(url, vendor, *KVStore)
	utils.ErrorCheck(err)
	return engine
}
