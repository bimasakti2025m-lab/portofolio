package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	suite.Suite
	router     *gin.Engine
	authUCMock *usecase.AuthenticateUsecaseMock
}

func (s *AuthControllerTestSuite) SetupTest() {
	s.router = gin.Default()
	s.authUCMock = new(usecase.AuthenticateUsecaseMock)
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (s *AuthControllerTestSuite) TestLoginHandler_Success() {
	payload := model.UserCredential{Username: "admin", Password: "admin"}
	expectedToken := "mock-jwt-token"
	s.authUCMock.On("Login", payload.Username, payload.Password).Return(expectedToken, nil)

	NewAuthController(s.authUCMock, s.router.Group("/api/v1")).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expectedToken, response["token"])
	s.authUCMock.AssertExpectations(s.T())
}

func (s *AuthControllerTestSuite) TestLoginHandler_InvalidCredentials() {
	payload := model.UserCredential{Username: "wrong", Password: "wrong"}
	expectedError := errors.New("invalid credentials")
	s.authUCMock.On("Login", payload.Username, payload.Password).Return("", expectedError)

	NewAuthController(s.authUCMock, s.router.Group("/api/v1")).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Contains(s.T(), response["err"], expectedError.Error())
	s.authUCMock.AssertExpectations(s.T())
}
