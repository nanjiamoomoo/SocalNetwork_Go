package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

//Router to direct request to different handler based on URL
func InitRouter() *mux.Router {
	router := mux.NewRouter()
	//Automatically maps the request URL to HTTP handler
	router.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("POST")
	return router
}