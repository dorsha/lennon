# lennon
Go utility for ElasticSearch and Bleve

##Status
- [x] ElasticSearch support
- [x] Bleve support

## Prerequisite
* ElasticSearch installed and running - *not relevant when using Bleve*

## Install
* go get -u -t -v github.com/dorsha/lennon
* go get ./...
* [optional] go install (from the root directory)

## Usage
```lennon --help ```  

**Index a given document**  
```lennon -vendor <elastic/bleve> -action index -url <search_engine_url> -document <path_to_document> ```

**Search in the indexed document**  
```lennon -vendor <elastic/bleve> -action search -url <search_engine_url> -query <search_query> ```  

##Examples (without running go install)
######ElasticSearch
**Index person.json file**  
```go run lennon.go -vendor elastic -action index -url http://192.168.1.26:9200 -document samples/persons.json ``` 

**Search for a person**  
```go run lennon.go -vendor elastic -action search -url http://192.168.1.26:9200 -query john ```

######Bleve
**Index person.json file**  
```go run lennon.go -vendor bleve -action index -document samples/persons.json ``` 

**Search for a person**  
```go run lennon.go -vendor bleve -action search -query john ```
