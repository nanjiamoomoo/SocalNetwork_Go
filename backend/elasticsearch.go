package backend

import (
    "context"
    "fmt"

    "socialnetwork_go/constants"
	"socialnetwork_go/util"

    "github.com/olivere/elastic/v7"
)

var (
	ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
	client *elastic.Client
}

//reference https://github.com/olivere/elastic/blob/release-branch.v7/example_test.go to see 
//example of how to create an index in ElasticSearch

//obtain a client and create indexes
func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
	//obtain a new client
	client, err := elastic.NewClient(elastic.SetURL(config.Address),
		elastic.SetBasicAuth(config.Username, config.Password), elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}

	//check if post index already exists in ElasticSearch
	exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	//create post index if it does not exist
	if !exists {		
		//for message:{"type":"text"} ES will be analized before storing into Inverted Index
		//this provides a quick keyword-based search
		mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" },
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
		
		_, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())

		if err != nil {
			panic(err)
		}
	}
	
	//check if user index already exists in ElasticSearch
	exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

	//create user index if it does not exist
	if !exists {
		mapping := `{
			"mappings": {
				"properties": {
					"username": {"type": "keyword"},
					"password": {"type": "keyword"},
					"age":      {"type": "long", "index": false},
					"gender":   {"type": "keyword", "index": false}
				}
			}
		}`
		
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("Indexes are created.")

	//ESbackend can be used to obtain the client object to do operations on the ElasticSearch
	ESBackend = &ElasticsearchBackend{client: client}
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	searchResult, err := backend.client.Search().
		Index(index). //specify the index to use for search
		Query(query). //sets the query to perform
		Pretty(true). //pretty print the formatted reponse JSON
		Do(context.Background()) //executes the search and returns a SearchResult.
	
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index().
        Index(index). //specify the index to save
        Id(id). //document ID
        BodyJson(i). 
        Do(context.Background())
    return err
}


func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
    _, err := backend.client.DeleteByQuery().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())

    return err
}