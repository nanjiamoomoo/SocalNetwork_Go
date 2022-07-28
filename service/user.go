package service

import (
	"fmt"
	"reflect"
	"socialnetwork_go/backend"
	"socialnetwork_go/constants"
	"socialnetwork_go/model"

	"github.com/olivere/elastic/v7"
)

//use username and passord to check if current user is a valid user
func CheckUser(username, password string) (bool, error) {
	//TermQuery finds documents that contain the exact term specified
	//NewTermQuery creates and initializes a new TermQuery.
	query := elastic.NewTermQuery("username", username)

	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	//check if the password in the token matches the password found in the index
	var utype model.User
	for _,item := range searchResult.Each(reflect.TypeOf(utype)) {
		u := item.(model.User)
		if u.Password == password {
			fmt.Printf("Login as %s\n", username)
			return true, nil
		}
	}
	return false, nil
}

//add a new user
//check if the username has been used. if yes, return; if not, create a new user.
func AddUser(user *model.User) (bool, error) {
	query := elastic.NewTermQuery("username", user.Username)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)

	if err != nil {
        return false, err
    }

	//if hits > 0, means the username has been used
    if searchResult.TotalHits() > 0 {
        return false, nil
    }

	//if hits = 0, add a new user in the "user" index
	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
        return false, err
    }
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil
}