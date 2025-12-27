package controller

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/usecase"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	suite.Suite
	authUseCase *usecase.AuthUseCaseMock
	controller  *AuthController
	router      *gin.Engine
}

func (suite *AuthControllerTestSuite) SetupTest() {
	suite.authUseCase = new(usecase.AuthUseCaseMock)
	suite.controller = &AuthController{
		authUc: suite.authUseCase,
	}

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	suite.router.POST("/register", suite.controller.registerHandler)
	suite.router.POST("/login", suite.controller.loginHandler)
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (suite *AuthControllerTestSuite) TestRegister_Success() {
	user := model.User{Username: "test", Email: "test@mail.com", Password: "password", Role: "user"}

	suite.authUseCase.On("Register", user.Username, user.Email, user.Password, user.Role).Return(user, nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *AuthControllerTestSuite) TestRegister_Conflict() {
	user := model.User{Username: "test", Email: "test@mail.com", Password: "password", Role: "user"}

	suite.authUseCase.On("Register", user.Username, user.Email, user.Password, user.Role).Return(model.User{}, errors.New("username 'test' is already taken"))

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusConflict, w.Code)
}

func (suite *AuthControllerTestSuite) TestLogin_Success() {
	user := model.User{Username: "test", Password: "password"}
	token := "mocked_token"

	suite.authUseCase.On("Login", user.Username, user.Password).Return(token, nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), token, response["token"])
}

func (suite *AuthControllerTestSuite) TestLogin_Failed() {
	user := model.User{Username: "test", Password: "wrongpassword"}

	suite.authUseCase.On("Login", user.Username, user.Password).Return("", errors.New("invalid credentials"))

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
