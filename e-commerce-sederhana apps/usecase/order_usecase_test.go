package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
	"E-commerce-Sederhana/utils/service/midtrans"
	"errors"
	"testing"

	"github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OrderUsecaseTestSuite struct {
	suite.Suite
	orderRepo       *repository.OrderRepositoryMock
	midtransService *midtrans.MidtransServiceMock
	usecase         OrderUsecase
}

func (suite *OrderUsecaseTestSuite) SetupTest() {
	suite.orderRepo = new(repository.OrderRepositoryMock)
	suite.midtransService = new(midtrans.MidtransServiceMock)
	suite.usecase = NewOrderUsecase(suite.orderRepo, suite.midtransService)
}

func TestOrderUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderUsecaseTestSuite))
}

func (suite *OrderUsecaseTestSuite) TestCreateOrder_Success() {
	order := model.Order{UserID: 1}
	createdOrder := order
	createdOrder.ID = 1
	createdOrder.TransactionIDMidtrans = "http://payment.url"

	suite.orderRepo.On("CreateOrder", mock.Anything).Return(createdOrder, nil)
	suite.midtransService.On("CreateTransaction", mock.Anything).Return(&snap.Response{RedirectURL: "http://payment.url"}, nil)

	result, err := suite.usecase.CreateOrder(order)
	suite.NoError(err)
	suite.Equal(createdOrder, result)
}

func (suite *OrderUsecaseTestSuite) TestCreateOrder_RepoError() {
	order := model.Order{UserID: 1}
	suite.orderRepo.On("CreateOrder", mock.Anything).Return(model.Order{}, errors.New("db error"))

	result, err := suite.usecase.CreateOrder(order)
	suite.Error(err)
	suite.Equal(model.Order{}, result)
}

func (suite *OrderUsecaseTestSuite) TestGetAllOrders_Success() {
	orders := []model.Order{{ID: 1, UserID: 1}}
	suite.orderRepo.On("GetAllOrders").Return(orders, nil)

	result, err := suite.usecase.GetAllOrders()
	suite.NoError(err)
	suite.Equal(orders, result)
}

func (suite *OrderUsecaseTestSuite) TestGetOrderById_Success() {
	order := model.Order{ID: 1, UserID: 1}
	suite.orderRepo.On("GetOrderById", 1).Return(order, nil)

	result, err := suite.usecase.GetOrderById(1)
	suite.NoError(err)
	suite.Equal(order, result)
}

func (suite *OrderUsecaseTestSuite) TestUpdateOrder_Success() {
	order := model.Order{ID: 1, StatusPesanan: "paid"}
	suite.orderRepo.On("UpdateOrder", order).Return(order, nil)

	result, err := suite.usecase.UpdateOrder(order)
	suite.NoError(err)
	suite.Equal(order, result)
}

func (suite *OrderUsecaseTestSuite) TestDeleteOrder_Success() {
	suite.orderRepo.On("DeleteOrder", 1).Return(nil)

	err := suite.usecase.DeleteOrder(1)
	suite.NoError(err)
}
