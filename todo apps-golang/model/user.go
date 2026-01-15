package model

type UserCredential struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}
