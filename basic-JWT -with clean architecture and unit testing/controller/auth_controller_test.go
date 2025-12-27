package controller

import (
	"basic-JWT/mock/controller_mock"
	"basic-JWT/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTest struct {
	suite.Suite
	authUc *controller_mock.AuthenticationUsecaseMock
	router *gin.Engine
	ac     *AuthController
}

func TestAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTest))
}

func (ac *AuthControllerTest) SetupTest() {
	ac.authUc = new(controller_mock.AuthenticationUsecaseMock)
	ac.router = gin.Default()
	rg := ac.router.Group("/api/v1")
	ac.ac = NewAuthController(rg, ac.authUc)
	ac.ac.Route()
}

func (ac *AuthControllerTest) TestRegisterHandler_Success() {
	user := model.User{Username: "testuser", Password: "testpassword"}
	ac.authUc.On("Register", user.Username, user.Password).Return(user, nil)

	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ac.router.ServeHTTP(w, req)

	ac.Equal(http.StatusOK, w.Code)
	ac.Contains(w.Body.String(), "success")
}

func (ac *AuthControllerTest) TestRegisterHandler_BadRequest() {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ac.router.ServeHTTP(w, req)

	ac.Equal(http.StatusBadRequest, w.Code)
	ac.Contains(w.Body.String(), "bad request")
}

func (ac *AuthControllerTest) TestRegisterHandler_UsernameTaken() {
	user := model.User{Username: "existinguser", Password: "testpassword"}
	ac.authUc.On("Register", user.Username, user.Password).Return(model.User{}, errors.New("username 'existinguser' is already taken"))

	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ac.router.ServeHTTP(w, req)

	ac.Equal(http.StatusConflict, w.Code)
	ac.Contains(w.Body.String(), "is already taken")
}

func (ac *AuthControllerTest) TestRegisterHandler_Failed() {
	user := model.User{Username: "testuser", Password: "testpassword"}
	ac.authUc.On("Register", user.Username, user.Password).Return(model.User{}, errors.New("some database error"))

	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ac.router.ServeHTTP(w, req)

	ac.Equal(http.StatusInternalServerError, w.Code)
	ac.Contains(w.Body.String(), "failed to register user")
}

func (ac *AuthControllerTest) TestLoginHandler_Success() {
	user := model.User{Username: "testuser", Password: "testpassword"}
	ac.authUc.On("Login", user.Username, user.Password).Return("testtoken", nil)

	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ac.router.ServeHTTP(w, req)

	ac.Equal(http.StatusOK, w.Code)
	ac.Contains(w.Body.String(), "testtoken")
}

func (ac *AuthControllerTest) TestLoginHandler_BadRequest() {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ac.router.ServeHTTP(w, req)

	ac.Equal(http.StatusBadRequest, w.Code)
	ac.Contains(w.Body.String(), "bad request")
}

func (ac *AuthControllerTest) TestLoginHandler_Failed() {
	user := model.User{Username: "testuser", Password: "testpassword"}
	ac.authUc.On("Login", user.Username, user.Password).Return("", errors.New("some database error"))

	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ac.router.ServeHTTP(w, req)

	ac.Equal(http.StatusInternalServerError, w.Code)
	ac.Contains(w.Body.String(), "failed to login user")
}
