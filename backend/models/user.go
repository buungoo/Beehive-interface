package models

// User is used to handle login.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
