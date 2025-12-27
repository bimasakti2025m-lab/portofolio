package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/entity/dto"
	"enigmacamp.com/livecode-catatan-keuangan/mock/usecase_mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTest struct {
	suite.Suite
	router *gin.Engine
	authUC *usecase_mock.AuthUsecaseMock
}

func (a *AuthControllerTest) SetupTest() {
	a.authUC = new(usecase_mock.AuthUsecaseMock)

	a.router = gin.Default()
	gin.SetMode(gin.TestMode)

	rg := a.router.Group("/api/v1")

	authC := NewAuthController(a.authUC, rg)
	a.router.POST("/api/v1/auth/login", authC.loginHandler) // Add this line to register the login handler
	a.router.POST("/api/v1/auth/register", authC.registerHandler)
}
func TestAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTest))
}

func (a *AuthControllerTest) TestRegisterHandler_Success() {
	payload := dto.AuthRequestDto{
		Username: "test-username",
		Password: "test-password",
	}

	a.authUC.On("Register", payload).Return(entity.User{
		Username: payload.Username,
		Password: payload.Password,
	}, nil)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	a.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/auth/register", &buf)
	a.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	record := httptest.NewRecorder()

	a.router.ServeHTTP(record, req)

	// Assertion to check status code
	a.Equal(http.StatusCreated, record.Code)
}

func (a *AuthControllerTest) TestRegisterHandler_Failed() {
	payload := dto.AuthRequestDto{
		Username: "test-username",
		Password: "test-password",
	}

	// Mock the Register method to return an error status bad request
	a.authUC.On("Register", payload).Return(entity.User{}, fmt.Errorf("failed"))

	// Convert payload to JSON

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	a.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/auth/register", &buf)
	a.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	record := httptest.NewRecorder()

	a.router.ServeHTTP(record, req)

	// Assertion to check status code
	a.Equal(http.StatusInternalServerError, record.Code)
}

func (a *AuthControllerTest) TestLoginHandler_Success() {
	// prepare 
	payload := dto.AuthRequestDto{
		Username: "test-username",
		Password: "test-password",
	}

	// mocking
	a.authUC.On("Login", payload).Return(dto.AuthResponseDto{
		Token: "test-token",
	}, nil)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	a.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", &buf)
	a.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	record := httptest.NewRecorder()

	a.router.ServeHTTP(record, req)

	// Assertion to check status code
	a.Equal(http.StatusCreated, record.Code)
}

func (a *AuthControllerTest) TestLoginHandler_Failed() {
	// prepare 
	payload := dto.AuthRequestDto{
		Username: "test-username",
		Password: "test-password",
	}

	// mocking
	a.authUC.On("Login", payload).Return(dto.AuthResponseDto{}, fmt.Errorf("failed"))

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	a.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", &buf)
	a.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	record := httptest.NewRecorder()

	a.router.ServeHTTP(record, req)

	// Assertion to check status code
	a.Equal(http.StatusInternalServerError, record.Code)
}

func (a *AuthControllerTest) TestLoginHandler_WrongPassword() {
	// prepare 
	payload := dto.AuthRequestDto{
		Username: "test-username",
		Password: "test-password",
	}

	// mocking
	a.authUC.On("Login", payload).Return(dto.AuthResponseDto{}, fmt.Errorf("wrong password"))

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	a.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", &buf)
	a.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	record := httptest.NewRecorder()

	a.router.ServeHTTP(record, req)

	// Assertion to check status code
	a.Equal(http.StatusInternalServerError, record.Code)
}

func (a *AuthControllerTest) TestLoginHandler_WrongUsername() {
	// prepare 
	payload := dto.AuthRequestDto{
		Username: "test-username",
		Password: "test-password",
	}

	// mocking
	a.authUC.On("Login", payload).Return(dto.AuthResponseDto{}, fmt.Errorf("wrong username"))

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	a.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", &buf)
	a.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	record := httptest.NewRecorder()

	a.router.ServeHTTP(record, req)

	// Assertion to check status code
	a.Equal(http.StatusInternalServerError, record.Code)
}

