package models

type User struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}
