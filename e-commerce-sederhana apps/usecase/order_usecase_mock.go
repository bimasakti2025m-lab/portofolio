package usecase

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type OrderUseCaseMock struct {
	mock.Mock
}

func (m *OrderUseCaseMock) GetAllOrders() ([]model.Order, error) {
	args := m.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *OrderUseCaseMock) GetOrderById(id int) (model.Order, error) {
	args := m.Called(id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *OrderUseCaseMock) CreateOrder(order model.Order) (model.Order, error) {
	args := m.Called(order)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *OrderUseCaseMock) UpdateOrder(order model.Order) (model.Order, error) {
	args := m.Called(order)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *OrderUseCaseMock) DeleteOrder(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
