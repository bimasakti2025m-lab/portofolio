package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OrderItemUsecaseTestSuite struct {
	suite.Suite
	orderItemRepo *repository.OrderItemRepositoryMock
	usecase       OrderItemUsecase
}

func (suite *OrderItemUsecaseTestSuite) SetupTest() {
	suite.orderItemRepo = new(repository.OrderItemRepositoryMock)
	suite.usecase = NewOrderItemUsecase(suite.orderItemRepo)
}

func TestOrderItemUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderItemUsecaseTestSuite))
}

func (suite *OrderItemUsecaseTestSuite) TestCreateOrderItem_Success() {
	orderItem := model.OrderItem{ProductID: 1, Quantity: 2}
	createdOrderItem := orderItem
	createdOrderItem.ID = 1

	suite.orderItemRepo.On("CreateOrderItem", &orderItem).Return(&createdOrderItem, nil)

	result, err := suite.usecase.CreateOrderItem(orderItem)
	suite.NoError(err)
	suite.Equal(&createdOrderItem, result)
}

func (suite *OrderItemUsecaseTestSuite) TestGetAllOrderItems_Success() {
	orderItems := []model.OrderItem{{ID: 1, ProductID: 1}}
	suite.orderItemRepo.On("GetAllOrderItems").Return(orderItems, nil)

	result, err := suite.usecase.GetAllOrderItems()
	suite.NoError(err)
	suite.Equal(orderItems, result)
}

func (suite *OrderItemUsecaseTestSuite) TestGetOrderItemByID_Success() {
	orderItem := model.OrderItem{ID: 1, ProductID: 1}
	suite.orderItemRepo.On("GetOrderItemByID", 1).Return(orderItem, nil)

	result, err := suite.usecase.GetOrderItemByID(1)
	suite.NoError(err)
	suite.Equal(orderItem, result)
}

func (suite *OrderItemUsecaseTestSuite) TestUpdateOrderItem_Success() {
	orderItem := model.OrderItem{ID: 1, Quantity: 5}
	updatedOrderItem := orderItem

	suite.orderItemRepo.On("UpdateOrderItem", &orderItem).Return(&updatedOrderItem, nil)

	result, err := suite.usecase.UpdateOrderItem(orderItem)
	suite.NoError(err)
	suite.Equal(&updatedOrderItem, result)
}

func (suite *OrderItemUsecaseTestSuite) TestDeleteOrderItem_Success() {
	suite.orderItemRepo.On("DeleteOrderItem", 1).Return(nil)

	err := suite.usecase.DeleteOrderItem(1)
	suite.NoError(err)
}

func (suite *OrderItemUsecaseTestSuite) TestDeleteOrderItem_Error() {
	suite.orderItemRepo.On("DeleteOrderItem", 1).Return(errors.New("failed"))

	err := suite.usecase.DeleteOrderItem(1)
	suite.Error(err)
}
