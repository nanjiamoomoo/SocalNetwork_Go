# SocialNetwork_Go

## **Project Description**
A social network web application that allows users to upload, search and delete posts (_This repository only includes the backend portion of the project_).

## **APIs**
APIs are defined under handler package
_`/signup`_ new user signup
_`/signin`_ sign in
_`/upload`_ allow users to upload post 
_`/search`_ allow users to search posts based on user and keywords
_`/post/{id}`_ allow users to delete post. This API only allows to delete the post id that belongs to the loggedin user

## **Packages**
**_`hanlder`_** 
Defines all the APIs handlder functions. Request data are preprocessed in the handler function and prepared to be passed into service layer. Handler handles the communication with frontend and all requests will go through this layer before sending to the service layer for processing. 

**_`service`_** 
Contains all the service logic and mediates communications between handler and backend. Queries are created in this layer and prepared to be passed to backend to access databases.

**_`backend`_** 
Communicates with ElasticSearch and GCS. All the CRUD methods are defined here to operate and return data in the database and GCS. 

**_`util, conf packages`_** Defines the configuration constants to be used and load the configuration when the server starts. 

**_`constants packages`_** Defines the name of the indices(databases) in ElasticSearch.

## **ElasticSearch and GCS**
**ElasticSearch** is used to store data posted by users. There are two indices created, "user" and "post".

**GCS** is used to store all media files posted by users. 

**_upload:_** When a post is uploaded, backend will generate a unique id and store the media file with the id to GCS. GCS will return the url info. Based on the id, user info, message input, url, and mediafile type, backend will wrap them up and store it as document in the post index.

**_search:_** When doing search, if the search is based on user, backend will search in post index and return all documents that matches the exact user; if the search is based on keywords, backend will return all documents that contain the keywords.

**_delete:_** When delete a post, backend will delele the doument based on the post id and the autenticated user's username. Only the document with the specified id and username can be removed from ElasticSearch.

## **Token-Based Authentication**
The project used JSON Web Token for server authentication. Token based authentication frees the server from saving a record of tokens. 

* A user enters their login credentials (=username + password).
* The server verifies the credentials are correct and created an encrypted and signed token with a private key ( { username: “abcd”, exp: “2021/1/1/10:00” }, private key => token).
* Client-side stores the token returned from the server.
* On subsequent requests, the token is decoded with the same private key and if valid the request is processed.
* Once a user logs out, the token is destroyed client-side, no interaction with the server is necessary.


## **Test**
Use Postman to test all APIs. 

**_`/signup`_** Create signup request by using Http POST Metthod. Send POST request with JSON format data in the request body to (http://Server_IP_ADDRESS:8080/signup).
 {
  "username": "admin",
  "password": "admin",
  "age": "30,
  "gender": "male"
  }

**_`/signin`_** Create signin request by using Http POST Metthod. Send POST request with JSON format data in the request body to (http://Server_IP_ADDRESS:8080/signin). A token string should be returned in response message.
 {
	"username":"admin",
	"password":"admin"
}

**_`/upload`_** Create upload request by using POST method and put http://Server_IP_ADDRESS:8080/upload in the address bar. In the Body section, choose form-data from the top and then add user, message, and media_file. For media_file, the type is “file” so that you may upload an image/video from your local disk. Under the Authorization tab, add a Bearer Token with the token string returned from signin.
 
**_`/search`_**  Create search request by using GET method, and put http://Server_IP_ADDRESS:8080/search?user=xxx in the address bar. Under the Authorization tab, add a Bearer Token with the token string returned from signin.

**_`/post/{id}`_** Create delete request by using DELETE method, and put http://Server_IP_ADDRESS:8080/post/put_postid_here in the address bar. Under the Authorization tab, add a Bearer Token with the token string returned from signin.

