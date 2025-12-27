package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"enigmacamp.com/url-shortener/middleware"
	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/usecase"
	modelutils "enigmacamp.com/url-shortener/utils/model_utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockAuthMiddleware digunakan untuk melewati otentikasi selama tes controller.
type MockAuthMiddlewareUrl struct{}

func (m *MockAuthMiddlewareUrl) RequireToken(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set claims palsu untuk mensimulasikan pengguna yang terotentikasi
		c.Set("claims", &modelutils.JwtPayloadClaims{UserId: 1, Role: "user"})
		c.Next()
	}
}

type UrlControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	usecaseMock *usecase.UrlUsecaseMock
	authMid     middleware.AuthMiddleware
}

func (s *UrlControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.usecaseMock = new(usecase.UrlUsecaseMock)
	s.authMid = &MockAuthMiddlewareUrl{} // Gunakan mock middleware
}

func TestUrlControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UrlControllerTestSuite))
}

func (s *UrlControllerTestSuite) TestCreateShortUrl_Success() {
	// Arrange
	payload := model.Url{LongUrl: "https://www.google.com"}
	expectedUrl := model.Url{
		Id:        1,
		LongUrl:   "https://www.google.com",
		ShortCode: "abc1234",
		UserId:    1, // User ID diambil dari mock claims
		CreatedAt: time.Now(),
	}

	// Matcher untuk memastikan UserId ditambahkan dengan benar di controller
	matcher := mock.MatchedBy(func(p model.Url) bool {
		return p.LongUrl == payload.LongUrl && p.UserId == expectedUrl.UserId
	})
	s.usecaseMock.On("CreateShortUrl", matcher).Return(expectedUrl, nil)

	NewUrlController(s.usecaseMock, s.router, s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	var actualUrl model.Url
	json.Unmarshal(w.Body.Bytes(), &actualUrl)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	assert.Equal(s.T(), expectedUrl.Id, actualUrl.Id)
	assert.Equal(s.T(), expectedUrl.ShortCode, actualUrl.ShortCode)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UrlControllerTestSuite) TestCreateShortUrl_BindingError() {
	// Arrange
	NewUrlController(s.usecaseMock, s.router, s.authMid).Route()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *UrlControllerTestSuite) TestCreateShortUrl_UseCaseError() {
	// Arrange
	payload := model.Url{LongUrl: "https://www.google.com"}
	s.usecaseMock.On("CreateShortUrl", mock.Anything).Return(model.Url{}, errors.New("usecase error"))

	NewUrlController(s.usecaseMock, s.router, s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Failed to create short URL")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UrlControllerTestSuite) TestRedirectToLongUrl_Success() {
	// Arrange
	shortCode := "abc1234"
	longUrl := "https://www.google.com"
	s.usecaseMock.On("GetLongUrl", shortCode).Return(longUrl, nil)

	NewUrlController(s.usecaseMock, s.router, s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/"+shortCode, nil)
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusMovedPermanently, w.Code)
	assert.Equal(s.T(), longUrl, w.Header().Get("Location"))
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UrlControllerTestSuite) TestRedirectToLongUrl_NotFound() {
	// Arrange
	shortCode := "notfound"
	expectedError := errors.New("url not found")
	s.usecaseMock.On("GetLongUrl", shortCode).Return("", expectedError)

	NewUrlController(s.usecaseMock, s.router, s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/"+shortCode, nil)
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusNotFound, w.Code)
	assert.Contains(s.T(), w.Body.String(), expectedError.Error())
	s.usecaseMock.AssertExpectations(s.T())
}
