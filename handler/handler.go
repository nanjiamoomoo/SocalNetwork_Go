package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"

)

//Router to direct request to different handler based on URL
func InitRouter() *mux.Router {


	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodES256,
	})

	router := mux.NewRouter()
	//Automatically maps the request URL to HTTP handler
	router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST") //protect upload with JWT
	router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET") //protect upload with JWT
	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
    router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")
	return router
}