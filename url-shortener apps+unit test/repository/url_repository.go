package repository

import (
	"database/sql"

	"enigmacamp.com/url-shortener/model"
)

type UrlRepository interface {
	Create(payload model.Url) (model.Url, error)
	GetByShortCode(shortCode string) (model.Url, error)
	IsShortCodeExist(shortCode string) (bool, error)
}

type urlRepository struct {
	db *sql.DB
}

func (r *urlRepository) Create(payload model.Url) (model.Url, error) {
	err := r.db.QueryRow("INSERT INTO mst_urls (long_url, short_code, user_id) VALUES ($1, $2, $3) RETURNING id, created_at", payload.LongUrl, payload.ShortCode, payload.UserId).Scan(&payload.Id, &payload.CreatedAt)
	return payload, err
}

func (r *urlRepository) GetByShortCode(shortCode string) (model.Url, error) {
	var url model.Url
	err := r.db.QueryRow("SELECT id, long_url, short_code, user_id, created_at FROM mst_urls WHERE short_code = $1", shortCode).Scan(&url.Id, &url.LongUrl, &url.ShortCode, &url.UserId, &url.CreatedAt)
	return url, err
}

func (r *urlRepository) IsShortCodeExist(shortCode string) (bool, error) {
	var exists bool // Menggunakan EXISTS untuk memeriksa keberadaan
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM mst_urls WHERE short_code = $1)", shortCode).Scan(&exists)
	return exists, err
}

func NewUrlRepository(db *sql.DB) UrlRepository {
	return &urlRepository{db: db}
}
