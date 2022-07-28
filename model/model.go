package model

// Define the model structs
// Field tags are a generic mechanism for encoding extra meta-data about a field
// JSON encoder uses this mechanism to let you override the key name in the generated JSON
type Post struct {
	Id      string `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Age      int64  `json:"age"`
    Gender   string `json:"gender"`
}

