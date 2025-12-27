package usecase

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"

	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/repository"
)

type UrlUsecase interface {
	CreateShortUrl(payload model.Url) (model.Url, error)
	GetLongUrl(shortCode string) (string, error)
}

type urlUsecase struct {
	repo repository.UrlRepository
}

const shortCodeLength = 7
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (u *urlUsecase) generateUniqueShortCode() (string, error) {
	for {
		b := make([]byte, shortCodeLength)
		for i := range b {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", err
			}
			b[i] = charset[num.Int64()]
		}
		shortCode := string(b)
		exists, err := u.repo.IsShortCodeExist(shortCode)
		if err != nil {
			return "", err
		}
		if !exists {
			return shortCode, nil
		}
	}
}

func (u *urlUsecase) CreateShortUrl(payload model.Url) (model.Url, error) {
	shortCode, err := u.generateUniqueShortCode()
	if err != nil {
		return model.Url{}, fmt.Errorf("failed to generate short code: %v", err)
	}
	payload.ShortCode = shortCode
	return u.repo.Create(payload)
}

func (u *urlUsecase) GetLongUrl(shortCode string) (string, error) {
	url, err := u.repo.GetByShortCode(shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("url with short code '%s' not found", shortCode)
		}
		return "", fmt.Errorf("failed to retrieve url: %v", err)
	}
	return url.LongUrl, nil
}

func NewUrlUsecase(repo repository.UrlRepository) UrlUsecase {
	return &urlUsecase{repo: repo}
}
