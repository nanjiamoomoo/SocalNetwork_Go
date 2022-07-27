package main

import (
	"fmt"
	"log"
	"net/http"

	"socialnetwork_go/handler"
	"socialnetwork_go/backend"
)

func main() {
	fmt.Println("started service!")

	//init client object for ElasticSearch and create indexes
	backend.InitElasticsearchBackend()
	//init client object for GCS
	backend.InitGCSBackend()

	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}