package usecase

import (
	"enigmacamp.com/url-shortener/model"
	"github.com/stretchr/testify/mock"
)

// UrlUsecaseMock adalah mock untuk UrlUsecase
type UrlUsecaseMock struct {
	mock.Mock
}

// CreateShortUrl adalah implementasi mock untuk metode CreateShortUrl
func (m *UrlUsecaseMock) CreateShortUrl(payload model.Url) (model.Url, error) {
	args := m.Called(payload)
	return args.Get(0).(model.Url), args.Error(1)
}

// GetLongUrl adalah implementasi mock untuk metode GetLongUrl
func (m *UrlUsecaseMock) GetLongUrl(shortCode string) (string, error) {
	args := m.Called(shortCode)
	return args.String(0), args.Error(1)
}
