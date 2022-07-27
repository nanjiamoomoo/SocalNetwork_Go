package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialnetwork_go/model"
	"socialnetwork_go/service"
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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
	w.Header().Set("Content-Type", "application/json")

	//get user from URL
	user := r.URL.Query().Get("user")
	//get keywords from URL
	keywords := r.URL.Query().Get("keywords")

	var posts []model.Post
	var err error

	//search based on either user or keywords
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeyWords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		fmt.Printf("Faied to read post from backend %v.\n", err)
		return
	}

	//convert Go object into JSON
	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
	}
	w.Write(js)
}