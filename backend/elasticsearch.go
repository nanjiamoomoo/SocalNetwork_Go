package backend

import (
    "context"
    "fmt"

    "socialnetwork_go/constants"

    "github.com/olivere/elastic/v7"
)

var (
	ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
	client *elastic.Client
}

//initialize the Elastic connection and create indexes
//reference https://github.com/olivere/elastic/blob/release-branch.v7/example_test.go to see 
//example of how to create an index in ElasticSearch
func InitElasticsearchBackend() {
	//obtain a new client
	client, err := elastic.NewClient(elastic.SetURL(constants.ES_URL),
		elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD))

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
		mapping := `
		{
			"mappings":{
				"properties":{
					"id":{
						"type":"keyword"
					},
					"user":{
						"type":"keyword"
					},
					"message":{
						"type":"text"
					},
					"url":{
						"type":"keyword",
						"index": false
					},
					"type":{
						"type":"keyword".
						"index": false
					}		
				}	
			}
		}
		`
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
        mapping := `
		{
			"mappings": {
				"properties":{
					"username":{
						"type":"keyword"
					},
					"password":{
						"type":"keyword"
					},
					"age":{
						"type":"long", 
						"index":false
					},
					"gender":{
						"type":"text", 
						"index":false
					}
				}
			}
        }
		`
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("Indexes are created.")

	//ESbackend can be used to obtain the client object to access the ElasticSearch
	ESBackend = &ElasticsearchBackend{client: client}
}