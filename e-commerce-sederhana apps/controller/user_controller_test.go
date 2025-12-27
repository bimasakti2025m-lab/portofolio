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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
	userUseCase *usecase.UserUseCaseMock
	controller  *UserController
	router      *gin.Engine
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.userUseCase = new(usecase.UserUseCaseMock)

	// Initialize controller with mocked usecase.
	// Dependencies not used in the handlers (rg, authMiddleware) can be nil for this unit test.
	suite.controller = &UserController{
		userUc: suite.userUseCase,
	}

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	// Register routes to test the handlers
	suite.router.POST("/users", suite.controller.createUserHandler)
	suite.router.GET("/users", suite.controller.getAllUsersHandler)
	suite.router.GET("/users/:username", suite.controller.getUserByUsernameHandler)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (suite *UserControllerTestSuite) TestCreateUser_Success() {
	user := model.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Role:     "user",
	}

	// Use mock.MatchedBy because the controller creates a new pointer for the user
	suite.userUseCase.On("Create", mock.MatchedBy(func(u *model.User) bool {
		return u.Username == user.Username && u.Email == user.Email
	})).Return(&user, nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])
}

func (suite *UserControllerTestSuite) TestCreateUser_BindError() {
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserControllerTestSuite) TestCreateUser_UseCaseError() {
	user := model.User{
		Username: "testuser",
		Password: "password",
	}

	suite.userUseCase.On("Create", mock.AnythingOfType("*model.User")).Return((*model.User)(nil), errors.New("failed to create user"))

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *UserControllerTestSuite) TestGetAllUsers_Success() {
	users := []model.User{
		{Username: "user1", Email: "user1@example.com"},
		{Username: "user2", Email: "user2@example.com"},
	}

	suite.userUseCase.On("GetAllUsers").Return(users, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string][]model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), users, response["users"])
}

func (suite *UserControllerTestSuite) TestGetAllUsers_Error() {
	suite.userUseCase.On("GetAllUsers").Return([]model.User(nil), errors.New("failed to get users"))

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *UserControllerTestSuite) TestGetUserByUsername_Success() {
	username := "testuser"
	user := model.User{Username: username, Email: "test@example.com"}

	suite.userUseCase.On("GetUserByUsername", username).Return(user, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users/"+username, nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user, response["user"])
}

func (suite *UserControllerTestSuite) TestGetUserByUsername_Error() {
	username := "unknown"
	suite.userUseCase.On("GetUserByUsername", username).Return(model.User{}, errors.New("failed to get user"))

	req, _ := http.NewRequest(http.MethodGet, "/users/"+username, nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
