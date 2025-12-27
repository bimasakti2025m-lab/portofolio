package controller

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/usecase"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OrderControllerTestSuite struct {
	suite.Suite
	orderUseCase *usecase.OrderUseCaseMock
	controller   *OrderController
	router       *gin.Engine
}

func (suite *OrderControllerTestSuite) SetupTest() {
	suite.orderUseCase = new(usecase.OrderUseCaseMock)
	suite.controller = &OrderController{
		orderUc: suite.orderUseCase,
	}

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	// Setup route khusus untuk create order agar bisa inject user ID ke context
	suite.router.POST("/orders", func(c *gin.Context) {
		c.Set("id", 1) // Simulasi user ID dari middleware auth
		suite.controller.createOrderHandler(c)
	})

	suite.router.GET("/orders", suite.controller.getAllOrdersHandler)
	suite.router.GET("/orders/:id", suite.controller.getOrderByIdHandler)
	suite.router.PUT("/orders/:id", suite.controller.updateOrderHandler)
	suite.router.DELETE("/orders/:id", suite.controller.deleteOrderHandler)
}

func TestOrderControllerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderControllerTestSuite))
}

func (suite *OrderControllerTestSuite) TestCreateOrder_Success() {
	order := model.Order{UserID: 10, Total: 50000}
	// UserID akan di-set oleh handler menjadi 1
	expectedOrderArg := order
	expectedOrderArg.UserID = 1

	createdOrder := expectedOrderArg
	createdOrder.ID = 100
	createdOrder.TransactionIDMidtrans = "http://payment.url"

	suite.orderUseCase.On("CreateOrder", expectedOrderArg).Return(createdOrder, nil)

	body, _ := json.Marshal(order)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "http://payment.url", response["payment_url"])
}

func (suite *OrderControllerTestSuite) TestGetAllOrders_Success() {
	orders := []model.Order{
		{ID: 1, UserID: 1, Total: 50000},
		{ID: 2, UserID: 1, Total: 75000},
	}

	suite.orderUseCase.On("GetAllOrders").Return(orders, nil)

	req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response []model.Order
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(suite.T(), response, 2)
}

func (suite *OrderControllerTestSuite) TestGetOrderById_Success() {
	order := model.Order{ID: 1, UserID: 1, Total: 50000}

	suite.orderUseCase.On("GetOrderById", 1).Return(order, nil)

	req, _ := http.NewRequest(http.MethodGet, "/orders/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *OrderControllerTestSuite) TestUpdateOrder_Success() {
	order := model.Order{StatusPesanan: "paid"}
	updatedOrder := order
	updatedOrder.ID = 1

	suite.orderUseCase.On("UpdateOrder", mock.MatchedBy(func(o model.Order) bool {
		return o.ID == 1 && o.StatusPesanan == "paid"
	})).Return(updatedOrder, nil)

	body, _ := json.Marshal(order)
	req, _ := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *OrderControllerTestSuite) TestDeleteOrder_Success() {
	suite.orderUseCase.On("DeleteOrder", 1).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/orders/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}
