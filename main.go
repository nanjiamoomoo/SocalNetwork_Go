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

	//inite client object and create indexes
	backend.InitElasticsearchBackend()

	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}