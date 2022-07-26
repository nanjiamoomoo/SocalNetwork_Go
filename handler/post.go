package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialnetwork_go/model"
)

//this method is to parse the body of the request to get a json object
//and convert the json object into the model object
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received on post request")
	decoder := json.NewDecoder(r.Body)

	var p model.Post
	if err := decoder.Decode(&p); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}