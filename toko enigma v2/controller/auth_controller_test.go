package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/toko-enigma/model"
	"enigmacamp.com/toko-enigma/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	usecaseMock *usecase.AuthenticateUsecaseMock
}

func (s *AuthControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.usecaseMock = new(usecase.AuthenticateUsecaseMock)
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (s *AuthControllerTestSuite) TestLoginHandler_Success() {
	payload := model.UserCredential{Username: "test", Password: "password"}
	expectedToken := "mocked-jwt-token"
	s.usecaseMock.On("Login", payload.Username, payload.Password).Return(expectedToken, nil)

	NewAuthController(s.usecaseMock, s.router.Group("/api/v1")).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expectedToken, response["token"])
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *AuthControllerTestSuite) TestLoginHandler_BindingError() {
	NewAuthController(s.usecaseMock, s.router.Group("/api/v1")).Route()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Binding JSON because")
}

func (s *AuthControllerTestSuite) TestLoginHandler_UsecaseError() {
	payload := model.UserCredential{Username: "test", Password: "password"}
	expectedError := "invalid username or password"
	s.usecaseMock.On("Login", payload.Username, payload.Password).Return("", errors.New(expectedError))

	NewAuthController(s.usecaseMock, s.router.Group("/api/v1")).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Login because")
	assert.Contains(s.T(), w.Body.String(), expectedError)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *AuthControllerTestSuite) TestRegisterHandler_Success() {
	payload := model.UserCredential{Username: "newuser", Password: "password", Role: "user"}
	expectedUser := model.UserCredential{Id: 1, Username: "newuser", Role: "user"}
	s.usecaseMock.On("Register", payload).Return(expectedUser, nil)

	NewAuthController(s.usecaseMock, s.router.Group("/api/v1")).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actualUser model.UserCredential
	json.Unmarshal(w.Body.Bytes(), &actualUser)

	assert.Equal(s.T(), http.StatusOK, w.Code) // Controller returns 200 on success
	assert.Equal(s.T(), expectedUser.Id, actualUser.Id)
	assert.Equal(s.T(), expectedUser.Username, actualUser.Username)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *AuthControllerTestSuite) TestRegisterHandler_BindingError() {
	NewAuthController(s.usecaseMock, s.router.Group("/api/v1")).Route()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Binding JSON because")
}

func (s *AuthControllerTestSuite) TestRegisterHandler_UsecaseError() {
	payload := model.UserCredential{Username: "newuser", Password: "password", Role: "user"}
	expectedError := "username already exists"
	s.usecaseMock.On("Register", payload).Return(model.UserCredential{}, errors.New(expectedError))

	NewAuthController(s.usecaseMock, s.router.Group("/api/v1")).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Register because")
	assert.Contains(s.T(), w.Body.String(), expectedError)
	s.usecaseMock.AssertExpectations(s.T())
}
