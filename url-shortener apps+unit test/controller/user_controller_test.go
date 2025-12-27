package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/url-shortener/middleware"
	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MockAuthMiddleware adalah mock untuk AuthMiddleware
type MockAuthMiddleware struct{}

func (m *MockAuthMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Langsung lanjutkan ke handler berikutnya tanpa validasi token
	}
}

type UserControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	usecaseMock *usecase.UserUseCaseMock
	authMid     middleware.AuthMiddleware
}

func (s *UserControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.usecaseMock = new(usecase.UserUseCaseMock)
	s.authMid = &MockAuthMiddleware{} // Gunakan mock middleware
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (s *UserControllerTestSuite) TestCreateUser_Success() {
	payload := model.UserCredential{Username: "test", Password: "password", Role: "user"}
	expectedUser := model.UserCredential{Id: 1, Username: "test", Role: "user"}
	s.usecaseMock.On("RegisterNewUser", payload).Return(expectedUser, nil)

	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actualUser model.UserCredential
	json.Unmarshal(w.Body.Bytes(), &actualUser)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	assert.Equal(s.T(), expectedUser.Id, actualUser.Id)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestCreateUser_Failed() {
	payload := model.UserCredential{Username: "test", Password: "password", Role: "user"}
	s.usecaseMock.On("RegisterNewUser", payload).Return(model.UserCredential{}, errors.New("some error"))

	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestCreateUser_InvalidJSON() {
	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *UserControllerTestSuite) TestCreateUser_BindingError() {
	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *UserControllerTestSuite) TestCreateUser_UseCaseError() {
	payload := model.UserCredential{Username: "test", Password: "password", Role: "user"}
	s.usecaseMock.On("RegisterNewUser", payload).Return(model.UserCredential{}, errors.New("some error"))

	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestGetAllUser_Success() {
	expectedUsers := []model.UserCredential{{Id: 1, Username: "test", Role: "user"}}
	s.usecaseMock.On("FindAllUser").Return(expectedUsers, nil)

	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actualUsers []model.UserCredential
	json.Unmarshal(w.Body.Bytes(), &actualUsers)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expectedUsers, actualUsers)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestGetAllUser_Empty() {
	s.usecaseMock.On("FindAllUser").Return([]model.UserCredential{}, nil)

	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Contains(s.T(), w.Body.String(), "List user empty")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestGetAllUser_UseCaseError() {
	s.usecaseMock.On("FindAllUser").Return(nil, errors.New("some error"))

	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestGetUserById_Success() {
	expectedUser := model.UserCredential{Id: 1, Username: "test", Role: "user"}
	s.usecaseMock.On("FindUserById", uint32(1)).Return(expectedUser, nil)

	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actualUser model.UserCredential
	json.Unmarshal(w.Body.Bytes(), &actualUser)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expectedUser, actualUser)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestGetUserById_UseCaseError() {
	s.usecaseMock.On("FindUserById", uint32(1)).Return(model.UserCredential{}, errors.New("some error"))

	NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestRoute() {
	// Hanya untuk memastikan rute terdaftar tanpa panic
	// Cakupan untuk fungsi Route() itu sendiri
	assert.NotPanics(s.T(), func() {
		NewUserController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()
	})
}
