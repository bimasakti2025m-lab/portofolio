package usecase

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type OrderItemUseCaseMock struct {
	mock.Mock
}

func (m *OrderItemUseCaseMock) CreateOrderItem(orderItem model.OrderItem) (*model.OrderItem, error) {
	args := m.Called(orderItem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderItem), args.Error(1)
}

func (m *OrderItemUseCaseMock) GetAllOrderItems() ([]model.OrderItem, error) {
	args := m.Called()
	return args.Get(0).([]model.OrderItem), args.Error(1)
}

func (m *OrderItemUseCaseMock) GetOrderItemByID(id int) (model.OrderItem, error) {
	args := m.Called(id)
	return args.Get(0).(model.OrderItem), args.Error(1)
}

func (m *OrderItemUseCaseMock) UpdateOrderItem(orderItem model.OrderItem) (*model.OrderItem, error) {
	args := m.Called(orderItem)
	return args.Get(0).(*model.OrderItem), args.Error(1)
}

func (m *OrderItemUseCaseMock) DeleteOrderItem(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
