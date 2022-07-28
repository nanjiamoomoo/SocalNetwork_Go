package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"socialnetwork_go/model"
	"socialnetwork_go/service"

	"github.com/form3tech-oss/jwt-go"
	"github.com/pborman/uuid"
)


var (
	mediaTypes = map[string]string{
		".jpeg": "image",
        ".jpg":  "image",
        ".gif":  "image",
        ".png":  "image",
        ".mov":  "video",
        ".mp4":  "video",
        ".avi":  "video",
        ".flv":  "video",
        ".wmv":  "video",
	}
)
//this method is to parse the body of the request to get a json object
//and convert the json object into the model object
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one upload request")

	//context() returns the request context
	user := r.Context().Value("user")
	fmt.Println(user)
	claims := user.(*jwt.Token).Claims //claims get the second segment of the token
	username := claims.(jwt.MapClaims)["username"]

	p := model.Post {
		Id: uuid.New(),
		User: username.(string),
		Message: r.FormValue("message"),
	}

	//FormFile returns the first file
	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename)
	if t, ok :=mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	err = service.SavePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to backend %v\n", err)
		return
	}

	fmt.Println("Post is saved successfully.")
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