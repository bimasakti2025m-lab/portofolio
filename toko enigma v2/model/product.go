package model

type Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Unit  string `json:"unit"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}
