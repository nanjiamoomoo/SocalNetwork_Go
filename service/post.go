package service

import (
	"mime/multipart"
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

func SavePost(post *model.Post, file multipart.File) error {
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
	if err != nil {
        return err
    }
    post.Url = medialink
	return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
}

func DeletePost(id string, user string) error {
    query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("id", id))
    query.Must(elastic.NewTermQuery("user", user))

    return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}