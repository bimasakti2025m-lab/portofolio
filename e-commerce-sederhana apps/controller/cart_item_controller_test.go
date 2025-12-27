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

type CartItemControllerTestSuite struct {
	suite.Suite
	cartItemUseCase *usecase.CartItemUseCaseMock
	controller      *CartItemController
	router          *gin.Engine
}

func (suite *CartItemControllerTestSuite) SetupTest() {
	suite.cartItemUseCase = new(usecase.CartItemUseCaseMock)
	suite.controller = &CartItemController{
		cartItemUc: suite.cartItemUseCase,
	}

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	suite.router.POST("/cart-items", suite.controller.createCartItemHandler)
	suite.router.GET("/cart-items", suite.controller.getAllCartItemsHandler)
	suite.router.GET("/cart-items/:id", suite.controller.getCartItemByIdHandler)
	suite.router.PUT("/cart-items/:id", suite.controller.updateCartItemHandler)
	suite.router.DELETE("/cart-items/:id", suite.controller.deleteCartItemHandler)
}

func TestCartItemControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CartItemControllerTestSuite))
}

func (suite *CartItemControllerTestSuite) TestCreateCartItem_Success() {
	cartItem := model.CartItem{ProductID: 1, Quantity: 2}
	createdCartItem := cartItem
	createdCartItem.ID = 1

	suite.cartItemUseCase.On("CreateCartItem", cartItem).Return(&createdCartItem, nil)

	body, _ := json.Marshal(cartItem)
	req, _ := http.NewRequest(http.MethodPost, "/cart-items", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *CartItemControllerTestSuite) TestGetAllCartItems_Success() {
	cartItems := []model.CartItem{{ID: 1, ProductID: 1}}
	suite.cartItemUseCase.On("GetAllCartItems").Return(cartItems, nil)

	req, _ := http.NewRequest(http.MethodGet, "/cart-items", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CartItemControllerTestSuite) TestGetCartItemById_Success() {
	cartItem := model.CartItem{ID: 1, ProductID: 1}
	suite.cartItemUseCase.On("GetCartItemByID", 1).Return(cartItem, nil)

	req, _ := http.NewRequest(http.MethodGet, "/cart-items/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CartItemControllerTestSuite) TestUpdateCartItem_Success() {
	cartItem := model.CartItem{Quantity: 5}
	updatedCartItem := cartItem
	updatedCartItem.ID = 1

	suite.cartItemUseCase.On("UpdateCartItem", mock.MatchedBy(func(c model.CartItem) bool {
		return c.ID == 1 && c.Quantity == 5
	})).Return(&updatedCartItem, nil)

	body, _ := json.Marshal(cartItem)
	req, _ := http.NewRequest(http.MethodPut, "/cart-items/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CartItemControllerTestSuite) TestDeleteCartItem_Success() {
	suite.cartItemUseCase.On("DeleteCartItem", 1).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/cart-items/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CartItemControllerTestSuite) TestDeleteCartItem_Error() {
	suite.cartItemUseCase.On("DeleteCartItem", 1).Return(errors.New("failed"))

	req, _ := http.NewRequest(http.MethodDelete, "/cart-items/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
