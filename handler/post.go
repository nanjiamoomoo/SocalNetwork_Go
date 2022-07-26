package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialnetwork_go/model"
)

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Receved on post request")
	decoder := json.NewDecoder(r.Body)

	var p model.Post
	var err: = decoder.Decode(&p);
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, "Post received: %s\n", p.Message)

}