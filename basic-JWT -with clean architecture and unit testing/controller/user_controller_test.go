package controller

import (
	"basic-JWT/middleware"
	"basic-JWT/mock/controller_mock"
	"basic-JWT/mock/service_mock"
	"basic-JWT/model"
	modelutils "basic-JWT/utils/model_utils"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserControllerTest struct {
	suite.Suite
	userUc         *controller_mock.UserUsecaseMock
	rg             *gin.Engine
	authMiddleware *middleware.AuthMiddleware
	jwtService     *service_mock.JWTServiceMock
	uc             *UserController
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTest))
}

func (uc *UserControllerTest) SetupTest() {
	uc.userUc = new(controller_mock.UserUsecaseMock)
	uc.rg = gin.Default()
	rg := uc.rg.Group("/api/v1")
	uc.jwtService = new(service_mock.JWTServiceMock)
	uc.authMiddleware = middleware.NewAuthMiddleware(uc.jwtService)
	uc.uc = NewUserController(rg, uc.userUc, uc.authMiddleware)
	uc.uc.Route() // Register routes
}

func (uc *UserControllerTest) TestCreateUserHandler_Success() {
	user := model.User{Username: "username", Password: "password", Role: "user"}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	// Mock the userUsecase.Create method
	uc.userUc.On("Create", &user).Return(&user, nil).Once()

	// Mock the jwtService.VerifyToken method
	uc.jwtService.On("VerifyToken", "dummy_admin_token").Return(&modelutils.JwtPayloadClaims{Role: "admin"}, nil).Once()

	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer dummy_admin_token") // Add a dummy token for authentication

	w := httptest.NewRecorder()
	uc.rg.ServeHTTP(w, req)

	uc.Equal(http.StatusOK, w.Code)
	uc.Contains(w.Body.String(), "success")
}

// masih error
// func (uc *UserControllerTest) TestCreateUserHandler_BadRequest() {
// 	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer([]byte("invalid-json")))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	uc.rg.ServeHTTP(w, req)

// 	uc.Equal(http.StatusBadRequest, w.Code)
// 	uc.Contains(w.Body.String(), "bad request")
// }

func (uc *UserControllerTest) TestCreateUserHandler_Failed() {
	user := model.User{Username: "username", Password: "password", Role: "user"}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	// Mock the userUsecase.Create method
	uc.userUc.On("Create", &user).Return(&user, errors.New("some database error"))

	// Mock the jwtService.VerifyToken method
	uc.jwtService.On("VerifyToken", "dummy_admin_token").Return(&modelutils.JwtPayloadClaims{Role: "admin"}, nil).Once()

	requestBody, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer dummy_admin_token") // Add a dummy token for authentication

	w := httptest.NewRecorder()
	uc.rg.ServeHTTP(w, req)

	uc.Equal(http.StatusInternalServerError, w.Code)
	uc.Contains(w.Body.String(), "failed to create user")
}

func (uc *UserControllerTest) TestGetAllUsersHandler_Success() {
	// Mock the userUsecase.GetAllUsers method
	uc.userUc.On("GetAllUsers").Return([]model.User{{Username: "user1"}, {Username: "user2"}}, nil)

	// Mock the jwtService.VerifyToken method
	uc.jwtService.On("VerifyToken", "dummy_admin_token").Return(&modelutils.JwtPayloadClaims{Role: "admin"}, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer dummy_admin_token") // Add a dummy token for authentication

	w := httptest.NewRecorder()
	uc.rg.ServeHTTP(w, req)

	uc.Equal(http.StatusOK, w.Code)
	uc.Contains(w.Body.String(), "user1")
	uc.Contains(w.Body.String(), "user2")
}

func (uc *UserControllerTest) TestGetAllUsersHandler_Failed() {
	user := []model.User{
		{
			Username: "user1",
			Password: "password1",
			Role:     "user",
		},
		{
			Username: "user2",
			Password: "password2",
			Role:     "user",
		},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user[0].Password), bcrypt.DefaultCost)
	user[0].Password = string(hashedPassword)
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user[1].Password), bcrypt.DefaultCost)
	user[1].Password = string(hashedPassword)

	// Mock the userUsecase.GetAllUsers method
	uc.userUc.On("GetAllUsers").Return(user, errors.New("some database error"))

	// Mock the jwtService.VerifyToken method
	uc.jwtService.On("VerifyToken", "dummy_admin_token").Return(&modelutils.JwtPayloadClaims{Role: "admin"}, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer dummy_admin_token") // Add a dummy token for authentication

	w := httptest.NewRecorder()
	uc.rg.ServeHTTP(w, req)

	uc.Equal(http.StatusInternalServerError, w.Code)
	uc.Contains(w.Body.String(), "failed to get users")

}

func (uc *UserControllerTest) TestGetUserByUsernameHandler_Success() {
	user := []model.User{
		{
			Username: "user1",
			Password: "user1password",
			Role:     "user",
		},
		{
			Username: "user2",
			Password: "user2password",
			Role:     "user",
		},
	}

	// Mock the userUsecase.GetUserByUsername method
	uc.userUc.On("GetUserByUsername", "user1").Return(user[0], nil)

	// Mock the jwtService.VerifyToken method
	uc.jwtService.On("VerifyToken", "dummy_admin_token").Return(&modelutils.JwtPayloadClaims{Role: "admin"}, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/user1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer dummy_admin_token") // Add a dummy token for authentication

	w := httptest.NewRecorder()
	uc.rg.ServeHTTP(w, req)

	uc.Equal(http.StatusOK, w.Code)
	uc.Contains(w.Body.String(), "user1")
	uc.Contains(w.Body.String(), "user1password")
}

func (uc *UserControllerTest) TestGetUserByUsernameHandler_Failed() {
	// Mock the userUsecase.GetUserByUsername method
	uc.userUc.On("GetUserByUsername", "user1").Return(model.User{}, errors.New("some database error"))

	// Mock the jwtService.VerifyToken method
	uc.jwtService.On("VerifyToken", "dummy_admin_token").Return(&modelutils.JwtPayloadClaims{Role: "admin"}, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/user1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer dummy_admin_token") // Add a dummy token for authentication

	w := httptest.NewRecorder()
	uc.rg.ServeHTTP(w, req)

	uc.Equal(http.StatusInternalServerError, w.Code)
	uc.Contains(w.Body.String(), "failed to get user")
}
