package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
    "time"

	"socialnetwork_go/model"
	"socialnetwork_go/service"

	"github.com/form3tech-oss/jwt-go"
)

var mySigningKey = []byte("secret")

/*
1. A user enters their login credentials (=username + password).
2. The server verifies the credentials are correct and created an encrypted and signed token with a private key ( { username: “abcd”, exp: “2021/1/1/10:00” }, private key => token).
3. Client-side stores the token returned from the server.
4. On subsequent requests, the token is decoded with the same private key and if valid the request is processed.
5. Once a user logs out, the token is destroyed client-side, no interaction with the server is necessary.
*/
func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
	w.Header().Set("Contect-Type", "text/plain")

	//decode JSON values from request body into user struct
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v\n", err)
        return
	}

	//check if current user is a valid user
	exists, err := service.CheckUser(user.Username, user.Password)
	if err != nil {
        http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
        return
    }

	if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User doesn't exists or wrong password\n")
        return
    }

	//created an encrypted and signed token with a private key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return
    }

	//return the private key
	w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")
    w.Header().Set("Content-Type", "text/plain")

	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v\n", err)
        return
	}

	if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

	//add a new user
	success, err := service.AddUser(&user)
	if err != nil {
        http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
        return
    }

	if !success {
        http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Println("User already exists")
        return
    }

	fmt.Printf("User added successfully: %s.\n", user.Username)
}