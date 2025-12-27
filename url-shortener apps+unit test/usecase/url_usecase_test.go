package usecase

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UrlUsecaseTestSuite struct {
	suite.Suite
	repoMock *repository.UrlRepositoryMock
	usecase  UrlUsecase
}

func (s *UrlUsecaseTestSuite) SetupTest() {
	s.repoMock = new(repository.UrlRepositoryMock)
	s.usecase = NewUrlUsecase(s.repoMock)
}

func TestUrlUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UrlUsecaseTestSuite))
}

func (s *UrlUsecaseTestSuite) TestCreateShortUrl_Success() {
	payload := model.Url{
		LongUrl: "https://www.google.com",
		UserId:  1,
	}

	// Kita tidak bisa memprediksi shortCode yang di-generate, jadi kita gunakan mock.AnythingOfType
	s.repoMock.On("IsShortCodeExist", mock.AnythingOfType("string")).Return(false, nil)

	// Ketika Create dipanggil, kita asumsikan berhasil dan mengembalikan data lengkap
	s.repoMock.On("Create", mock.MatchedBy(func(p model.Url) bool {
		return p.LongUrl == payload.LongUrl && p.UserId == payload.UserId
	})).Return(model.Url{
		Id:        1,
		LongUrl:   payload.LongUrl,
		ShortCode: "abcdefg", // Contoh shortcode
		UserId:    payload.UserId,
		CreatedAt: time.Now(),
	}, nil)

	createdUrl, err := s.usecase.CreateShortUrl(payload)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), createdUrl)
	assert.Equal(s.T(), payload.LongUrl, createdUrl.LongUrl)
	assert.NotEmpty(s.T(), createdUrl.ShortCode)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UrlUsecaseTestSuite) TestCreateShortUrl_FailGenerateCode() {
	payload := model.Url{
		LongUrl: "https://www.google.com",
		UserId:  1,
	}

	// Simulasikan error saat memeriksa keberadaan short code
	expectedError := errors.New("database error")
	s.repoMock.On("IsShortCodeExist", mock.AnythingOfType("string")).Return(false, expectedError)

	createdUrl, err := s.usecase.CreateShortUrl(payload)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to generate short code")
	assert.Equal(s.T(), model.Url{}, createdUrl) // Pastikan objek URL kosong
}

func (s *UrlUsecaseTestSuite) TestGetLongUrl_Success() {
	shortCode := "xyz123"
	expectedUrl := model.Url{
		Id:        1,
		LongUrl:   "https://www.enigmacamp.com",
		ShortCode: shortCode,
		UserId:    1,
	}
	s.repoMock.On("GetByShortCode", shortCode).Return(expectedUrl, nil)

	longUrl, err := s.usecase.GetLongUrl(shortCode)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUrl.LongUrl, longUrl)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UrlUsecaseTestSuite) TestGetLongUrl_NotFound() {
	shortCode := "notfound"
	s.repoMock.On("GetByShortCode", shortCode).Return(model.Url{}, sql.ErrNoRows)

	longUrl, err := s.usecase.GetLongUrl(shortCode)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "not found")
	assert.Empty(s.T(), longUrl)
	s.repoMock.AssertExpectations(s.T())
}