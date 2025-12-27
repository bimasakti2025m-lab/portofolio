package repository

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) GetAllOrders() ([]model.Order, error) {
	args := m.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetOrderById(id int) (model.Order, error) {
	args := m.Called(id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *OrderRepositoryMock) CreateOrder(order model.Order) (model.Order, error) {
	args := m.Called(order)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *OrderRepositoryMock) UpdateOrder(order model.Order) (model.Order, error) {
	args := m.Called(order)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *OrderRepositoryMock) DeleteOrder(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
