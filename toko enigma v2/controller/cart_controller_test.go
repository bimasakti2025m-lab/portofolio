package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/toko-enigma/middleware"
	"enigmacamp.com/toko-enigma/model"
	"enigmacamp.com/toko-enigma/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CartControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	usecaseMock *usecase.CartUseCaseMock
	authMid     middleware.AuthMiddleware
}

func (s *CartControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.usecaseMock = new(usecase.CartUseCaseMock)
	s.authMid = &MockAuthMiddleware{}
}

func TestCartControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CartControllerTestSuite))
}

func (s *CartControllerTestSuite) TestCreateCart_Success() {
	payload := model.Cart{UserId: "user-1", TotalPrice: 15000}
	expectedCart := model.Cart{Id: 1, UserId: "user-1", TotalPrice: 15000}
	s.usecaseMock.On("CreateNewCart", payload).Return(expectedCart, nil)

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/carts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actualCart model.Cart
	json.Unmarshal(w.Body.Bytes(), &actualCart)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	assert.Equal(s.T(), expectedCart, actualCart)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestCreateCart_BindingError() {
	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/carts", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Failed to bind JSON")
}

func (s *CartControllerTestSuite) TestCreateCart_UsecaseError() {
	payload := model.Cart{UserId: "user-1"}
	s.usecaseMock.On("CreateNewCart", payload).Return(model.Cart{}, errors.New("some error"))

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/carts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Failed to create cart")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestGetAllCarts_Success() {
	expectedCarts := []model.Cart{{Id: 1, UserId: "user-1"}}
	s.usecaseMock.On("GetAllCart").Return(expectedCarts, nil)

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/carts", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actualCarts []model.Cart
	json.Unmarshal(w.Body.Bytes(), &actualCarts)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expectedCarts, actualCarts)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestGetAllCarts_Empty() {
	s.usecaseMock.On("GetAllCart").Return([]model.Cart{}, nil)

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/carts", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Contains(s.T(), w.Body.String(), "List cart empty")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestGetAllCarts_UsecaseError() {
	s.usecaseMock.On("GetAllCart").Return(nil, errors.New("some error"))

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/carts", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Failed to retrieve data cart")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestGetCartById_Success() {
	expectedCart := model.Cart{Id: 1, UserId: "user-1"}
	s.usecaseMock.On("GetCartById", 1).Return(expectedCart, nil)

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/carts/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actualCart model.Cart
	json.Unmarshal(w.Body.Bytes(), &actualCart)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expectedCart, actualCart)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestGetCartById_UsecaseError() {
	s.usecaseMock.On("GetCartById", 1).Return(model.Cart{}, errors.New("some error"))

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/carts/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Failed to get cart by ID")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestUpdateCart_Success() {
	payload := model.Cart{Id: 1, UserId: "user-updated"}
	s.usecaseMock.On("UpdateCartById", payload).Return(payload, nil)

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/carts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actualCart model.Cart
	json.Unmarshal(w.Body.Bytes(), &actualCart)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), payload, actualCart)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestUpdateCart_UsecaseError() {
	payload := model.Cart{Id: 1, UserId: "user-updated"}
	expectedError := "Cart with ID : 1 not found."
	s.usecaseMock.On("UpdateCartById", payload).Return(model.Cart{}, errors.New(expectedError))

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/carts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Failed to update cart")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestDeleteCart_Success() {
	s.usecaseMock.On("DeleteCartById", 1).Return(nil)

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/carts/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), `"Cart deleted successfully."`, w.Body.String())
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *CartControllerTestSuite) TestDeleteCart_UsecaseError() {
	s.usecaseMock.On("DeleteCartById", 1).Return(errors.New("some error"))

	NewCartController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/carts/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Failed to delete todo")
	s.usecaseMock.AssertExpectations(s.T())
}
