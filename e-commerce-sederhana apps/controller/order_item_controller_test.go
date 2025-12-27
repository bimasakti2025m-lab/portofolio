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

type OrderItemControllerTestSuite struct {
	suite.Suite
	orderItemUseCase *usecase.OrderItemUseCaseMock
	controller       *OrderItemController
	router           *gin.Engine
}

func (suite *OrderItemControllerTestSuite) SetupTest() {
	suite.orderItemUseCase = new(usecase.OrderItemUseCaseMock)
	suite.controller = &OrderItemController{
		orderItemUc: suite.orderItemUseCase,
	}

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	suite.router.POST("/order-items", suite.controller.createOrderItemHandler)
	suite.router.GET("/order-items", suite.controller.getAllOrderItemsHandler)
	suite.router.GET("/order-items/:id", suite.controller.getOrderItemByIdHandler)
	suite.router.PUT("/order-items/:id", suite.controller.updateOrderItemHandler)
	suite.router.DELETE("/order-items/:id", suite.controller.deleteOrderItemHandler)
}

func TestOrderItemControllerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderItemControllerTestSuite))
}

func (suite *OrderItemControllerTestSuite) TestCreateOrderItem_Success() {
	orderItem := model.OrderItem{ProductID: 1, Quantity: 2}
	createdOrderItem := orderItem
	createdOrderItem.ID = 1

	suite.orderItemUseCase.On("CreateOrderItem", orderItem).Return(&createdOrderItem, nil)

	body, _ := json.Marshal(orderItem)
	req, _ := http.NewRequest(http.MethodPost, "/order-items", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *OrderItemControllerTestSuite) TestGetAllOrderItems_Success() {
	orderItems := []model.OrderItem{{ID: 1, ProductID: 1}}
	suite.orderItemUseCase.On("GetAllOrderItems").Return(orderItems, nil)

	req, _ := http.NewRequest(http.MethodGet, "/order-items", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *OrderItemControllerTestSuite) TestGetOrderItemById_Success() {
	orderItem := model.OrderItem{ID: 1, ProductID: 1}
	suite.orderItemUseCase.On("GetOrderItemByID", 1).Return(orderItem, nil)

	req, _ := http.NewRequest(http.MethodGet, "/order-items/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *OrderItemControllerTestSuite) TestUpdateOrderItem_Success() {
	// Note: Controller saat ini tidak mengambil ID dari parameter URL, melainkan dari body JSON
	orderItem := model.OrderItem{ID: 1, Quantity: 5}
	updatedOrderItem := orderItem

	suite.orderItemUseCase.On("UpdateOrderItem", orderItem).Return(&updatedOrderItem, nil)

	body, _ := json.Marshal(orderItem)
	req, _ := http.NewRequest(http.MethodPut, "/order-items/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *OrderItemControllerTestSuite) TestDeleteOrderItem_Success() {
	suite.orderItemUseCase.On("DeleteOrderItem", 1).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/order-items/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *OrderItemControllerTestSuite) TestDeleteOrderItem_Error() {
	suite.orderItemUseCase.On("DeleteOrderItem", 1).Return(errors.New("failed"))

	req, _ := http.NewRequest(http.MethodDelete, "/order-items/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
