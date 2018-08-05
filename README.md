# lennon
Go utility for ElasticSearch and Bleve

## Status
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
<sub>```lennon --help ```  </sub>

**Index a given document**  
<sub>```lennon -vendor <elastic/bleve> -action index -document <path_to_document> [-url <es_engine_url>] [-store <boltdb/goleveldb>]```</sub>

**Index documents inside folder**  
<sub>```lennon -vendor <elastic/bleve> -action index -folder <path_to_folder> [-url <es_engine_url>] [-store <boltdb/goleveldb>] ```</sub>

**Search in the indexed document**  
<sub>```lennon -vendor <elastic/bleve> -action search -query <search_query> [-url <es_engine_url>] [-store <boltdb/goleveldb>] ```</sub>  

**Delete index**  
<sub>```lennon -vendor <elastic/bleve> -action deleteIndex [-url <es_engine_url>] [-store <boltdb/goleveldb>] ```</sub>  

## Examples (without running go install)
### ElasticSearch
**Index person.json file**  
<sub>```go run lennon.go -vendor elastic -action index -url http://192.168.1.26:9200 -document samples/persons.json ```</sub> 

**Search for a person**  
<sub>```go run lennon.go -vendor elastic -action search -url http://192.168.1.26:9200 -query john ```</sub>

**Delete the index**  
<sub>```go run lennon.go -vendor elastic -action deleteIndex -url http://192.168.1.18:9200 ```</sub>

### Bleve
**Index person.json file (LevelDB)**  
<sub>```go run lennon.go -vendor bleve -action index -document samples/persons.json ```</sub>

**Index person.json file (BoltDB)**  
<sub>```go run lennon.go -vendor bleve -action index -document samples/persons.json -store boltdb ```</sub> 

**Search for a person (LevelDB)**  
<sub>```go run lennon.go -vendor bleve -action search -query john ```</sub>

**Delete the index (LevelDB)**  
<sub>```go run lennon.go -vendor bleve -action deleteIndex ```</sub>
