# lennon
Go utility for ElasticSearch and Bleve

##Status
- [x] ElasticSearch support
- [x] Bleve support
- [x] Index a document file
- [x] Index a folder with documents
- [x] Delete indexes
- [x] Search

## Prerequisite
* ElasticSearch installed and running - *not relevant when using Bleve*

## Install
* go get -u -t -v github.com/dorsha/lennon
* go get ./...
* [optional] go install (from the root directory)

## Usage
```lennon --help ```  

**Index a given document**  
```lennon -vendor <elastic/bleve> -action index [-url <es_engine_url>] -document <path_to_document> ```

**Index documents inside folder**  
```lennon -vendor <elastic/bleve> -action index [-url <es_engine_url>] -folder <path_to_folder> ```

**Search in the indexed document**  
```lennon -vendor <elastic/bleve> -action search [-url <es_engine_url>] -query <search_query> ```  

**Delete index**  
```lennon -vendor <elastic/bleve> -action deleteIndex [-url <es_engine_url>] ```  

##Examples (without running go install)
######ElasticSearch
**Index person.json file**  
```go run lennon.go -vendor elastic -action index -url http://192.168.1.26:9200 -document samples/persons.json ``` 

**Search for a person**  
```go run lennon.go -vendor elastic -action search -url http://192.168.1.26:9200 -query john ```

**Delete the index**  
```go run lennon.go -vendor elastic -action deleteIndex -url http://192.168.1.18:9200 ```

######Bleve
**Index person.json file**  
```go run lennon.go -vendor bleve -action index -document samples/persons.json ``` 

**Search for a person**  
```go run lennon.go -vendor bleve -action search -query john ```

**Delete the index**  
```go run lennon.go -vendor bleve -action deleteIndex ```
