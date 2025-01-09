package models

// User and chat models

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	Id   int64  `json:"id"`
	Type string `json:"type"` // private, group, supergroup, channel
}
