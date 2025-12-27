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

type CartControllerTestSuite struct {
	suite.Suite
	cartUseCase *usecase.CartUseCaseMock
	controller  *CartController
	router      *gin.Engine
}

func (suite *CartControllerTestSuite) SetupTest() {
	suite.cartUseCase = new(usecase.CartUseCaseMock)
	suite.controller = &CartController{
		cartUc: suite.cartUseCase,
	}

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	suite.router.POST("/carts", suite.controller.createCartHandler)
	suite.router.GET("/carts", suite.controller.getAllCartsHandler)
	suite.router.GET("/carts/:id", suite.controller.getCartByIdHandler)
	suite.router.PUT("/carts/:id", suite.controller.updateCartHandler)
	suite.router.DELETE("/carts/:id", suite.controller.deleteCartHandler)
}

func TestCartControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CartControllerTestSuite))
}

func (suite *CartControllerTestSuite) TestCreateCart_Success() {
	cart := model.Cart{UserID: 1}
	createdCart := cart
	createdCart.ID = 1

	suite.cartUseCase.On("CreateCart", cart).Return(&createdCart, nil)

	body, _ := json.Marshal(cart)
	req, _ := http.NewRequest(http.MethodPost, "/carts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *CartControllerTestSuite) TestGetAllCarts_Success() {
	carts := []model.Cart{{ID: 1, UserID: 1}}
	suite.cartUseCase.On("GetAllCarts").Return(carts, nil)

	req, _ := http.NewRequest(http.MethodGet, "/carts", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CartControllerTestSuite) TestGetCartById_Success() {
	cart := model.Cart{ID: 1, UserID: 1}
	suite.cartUseCase.On("GetCartByID", 1).Return(cart, nil)

	req, _ := http.NewRequest(http.MethodGet, "/carts/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CartControllerTestSuite) TestUpdateCart_Success() {
	cart := model.Cart{UserID: 2}
	updatedCart := cart
	updatedCart.ID = 1

	suite.cartUseCase.On("UpdateCart", mock.MatchedBy(func(c model.Cart) bool {
		return c.ID == 1 && c.UserID == 2
	})).Return(&updatedCart, nil)

	body, _ := json.Marshal(cart)
	req, _ := http.NewRequest(http.MethodPut, "/carts/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CartControllerTestSuite) TestDeleteCart_Success() {
	suite.cartUseCase.On("DeleteCart", 1).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/carts/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CartControllerTestSuite) TestDeleteCart_Error() {
	suite.cartUseCase.On("DeleteCart", 1).Return(errors.New("failed"))

	req, _ := http.NewRequest(http.MethodDelete, "/carts/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
