package repository

import (
	"enigmacamp.com/url-shortener/model"
	"github.com/stretchr/testify/mock"
)

// UrlRepositoryMock adalah mock untuk UrlRepository
type UrlRepositoryMock struct {
	mock.Mock
}

// Create adalah implementasi mock untuk metode Create
func (m *UrlRepositoryMock) Create(payload model.Url) (model.Url, error) {
	args := m.Called(payload)
	return args.Get(0).(model.Url), args.Error(1)
}

// GetByShortCode adalah implementasi mock untuk metode GetByShortCode
func (m *UrlRepositoryMock) GetByShortCode(shortCode string) (model.Url, error) {
	args := m.Called(shortCode)
	return args.Get(0).(model.Url), args.Error(1)
}

// IsShortCodeExist adalah implementasi mock untuk metode IsShortCodeExist
func (m *UrlRepositoryMock) IsShortCodeExist(shortCode string) (bool, error) {
	args := m.Called(shortCode)
	return args.Bool(0), args.Error(1)
}
