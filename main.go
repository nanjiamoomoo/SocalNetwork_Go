package main

import (
	"fmt"
	"log"
	"net/http"

	"socialnetwork_go/util"
	"socialnetwork_go/handler"
	"socialnetwork_go/backend"
)

func main() {
	fmt.Println("started service!")

	config, err := util.LoadApplicationConfig("conf", "deploy.yml")
    if err != nil {
        panic(err)
    }


	//init client object for ElasticSearch and create indexes
	backend.InitElasticsearchBackend(config.ElasticsearchConfig)
	//init client object for GCS
	backend.InitGCSBackend(config.GCSConfig)

	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter(config.TokenConfig)))
}