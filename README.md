# lennon
Go utility for ElasticSearch and Bleve

## Prerequisite
* ElasticSearch installed and running

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
**Index person.json file (ElasticSearch)**
```go run lennon.go -vendor elastic -action index -url http://192.168.1.26:9200 -document persons.json ``` 

**Search for a person (ElasticSearch)**
```go run lennon.go -vendor elastic -action search -url http://192.168.1.26:9200 -query jhon ```
