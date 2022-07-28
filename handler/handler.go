package handler

import (
	"net/http"

	"socialnetwork_go/util"

	"github.com/gorilla/mux"
	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"

)

var mySigningKey []byte

//Router to direct request to different handler based on URL
func InitRouter(config *util.TokenInfo) http.Handler{
	mySigningKey = []byte(config.Secret)
	
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter()
	//Automatically maps the request URL to HTTP handler
	router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST") //protect upload with JWT
	router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET") //protect upload with JWT
	router.Handle("/post/{id}", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE") //protect delete with JWT

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
    router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

	//make the backend support CORS requests
	originsOk := handlers.AllowedOrigins([]string{"*"})//support all
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"}) //methods support

	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}