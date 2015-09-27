# lennon
Go utility for ElasticSearch and Bleve

##Status
- [x] ElasticSearch support
- [x] Bleve support
- [x] Index a document file
- [x] Index a folder with documents
- [x] Delete indexes
- [x] Search
- [x] Batch support
- [x] KV Store support for Bleve (LevelDB and BoltDB)

## Prerequisite
* ElasticSearch installed and running - *not relevant when using Bleve*

## Install
* go get -u -t -v github.com/dorsha/lennon
* go get ./...
* [optional] go install (from the root directory)

## Usage
```lennon --help ```  

**Index a given document**  
```lennon -vendor <elastic/bleve> -action index -document <path_to_document> [-url <es_engine_url>] [-store <boltdb/goleveldb>]```

**Index documents inside folder**  
```lennon -vendor <elastic/bleve> -action index -folder <path_to_folder> [-url <es_engine_url>] [-store <boltdb/goleveldb>] ```

**Search in the indexed document**  
```lennon -vendor <elastic/bleve> -action search -query <search_query> [-url <es_engine_url>] [-store <boltdb/goleveldb>] ```  

**Delete index**  
```lennon -vendor <elastic/bleve> -action deleteIndex [-url <es_engine_url>] [-store <boltdb/goleveldb>] ```  

##Examples (without running go install)
######ElasticSearch
**Index person.json file**  
```go run lennon.go -vendor elastic -action index -url http://192.168.1.26:9200 -document samples/persons.json ``` 

**Search for a person**  
```go run lennon.go -vendor elastic -action search -url http://192.168.1.26:9200 -query john ```

**Delete the index**  
```go run lennon.go -vendor elastic -action deleteIndex -url http://192.168.1.18:9200 ```

######Bleve
**Index person.json file (LevelDB)**  
```go run lennon.go -vendor bleve -action index -document samples/persons.json ``` 

**Index person.json file (BoltDB)**  
```go run lennon.go -vendor bleve -action index -document samples/persons.json -store boltdb ``` 

**Search for a person (LevelDB)**  
```go run lennon.go -vendor bleve -action search -query john ```

**Delete the index (LevelDB)**  
```go run lennon.go -vendor bleve -action deleteIndex ```
