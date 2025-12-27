package model

type Book struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ReleaseYear int    `json:"releaseYear"`
	Pages       int    `json:"pages"`
}
