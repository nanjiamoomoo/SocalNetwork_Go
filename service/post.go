package service

import (
	"fmt"
	"reflect"
	"socialnetwork_go/backend"
	"socialnetwork_go/constants"
	"socialnetwork_go/model"

	"github.com/olivere/elastic/v7"
)

//search posts posted by user based on the "user" property
func SearchPostsByUser(user string) ([]model.Post, error) {
	//createa a new query based on the "user" property
	query := elastic.NewTermQuery("user", user)

	searchResult, err:= backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func SearchPostsByKeyWords(keywords string) ([]model.Post, error) {
	//createa a keywords-based query based on the "message" property
	query := elastic.NewMatchQuery("message", keywords)
	query.Operator("AND")
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
    // and all kinds of other information from Elasticsearch.
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil

}

//convert search results from elasticsearch to model.Post
func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
	var ptype model.Post
	var posts []model.Post

	//each iterates over hits in a search result.
	//reflect.TypeOf() get the reflection type
	for _,item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(model.Post)
		posts = append(posts, p)
	}
	return posts
}