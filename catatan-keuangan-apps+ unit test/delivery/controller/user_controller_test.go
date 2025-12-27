package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/livecode-catatan-keuangan/delivery/middleware"
	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/mock/usecase_mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type UserControllerTest struct {
	suite.Suite
	router *gin.Engine
	userUC *usecase_mock.UserUsecaseMock
	am     *middleware.AuthMiddleware
}

func (u *UserControllerTest) SetupTest() {
	u.userUC = new(usecase_mock.UserUsecaseMock)
	u.am = new(middleware.AuthMiddleware)

	u.router = gin.Default()
	gin.SetMode(gin.TestMode)

	rg := u.router.Group("/api/v1")

	userC := NewUserController(u.userUC, rg, *u.am)
	rg.GET("/users/:id", userC.getHandler) // Corrected: Register handler on the router group
}
func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTest))
}

func (u *UserControllerTest) TestGetHandler_Success() {
	// Prepare
	u.userUC.On("FindUserByID", "uuid-user-test").Return(entity.User{
		ID:       "uuid-user-test",
		Username: "success",
		Password: "password success",
	}, nil).Once()

	req, err := http.NewRequest("GET", "/api/v1/users/uuid-user-test", nil)
	u.NoError(err)

	record := httptest.NewRecorder()

	u.router.ServeHTTP(record, req)

	// Assertion to check status code
	u.Equal(http.StatusOK, record.Code)
	}

func (u *UserControllerTest) TestGetHandler_Failed() {
	// Prepare
	u.userUC.On("FindUserByID", "uuid-user-failed").Return(entity.User{}, fmt.Errorf("failed")).Once()

	req, err := http.NewRequest("GET", "/api/v1/users/uuid-user-failed", nil)
	u.NoError(err)

	record := httptest.NewRecorder()

	u.router.ServeHTTP(record, req)

	// Assertion to check status code
	u.Equal(http.StatusNotFound, record.Code)
}

