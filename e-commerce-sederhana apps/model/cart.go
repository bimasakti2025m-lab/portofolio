package model

type Cart struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	Status    string `json:"status"`
}