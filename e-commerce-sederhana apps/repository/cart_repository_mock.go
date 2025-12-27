package repository

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type CartRepositoryMock struct {
	mock.Mock
}

func (m *CartRepositoryMock) CreateCart(cart *model.Cart) (*model.Cart, error) {
	args := m.Called(cart)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *CartRepositoryMock) GetAllCarts() ([]model.Cart, error) {
	args := m.Called()
	return args.Get(0).([]model.Cart), args.Error(1)
}

func (m *CartRepositoryMock) GetCartByID(id int) (model.Cart, error) {
	args := m.Called(id)
	return args.Get(0).(model.Cart), args.Error(1)
}

func (m *CartRepositoryMock) UpdateCart(cart *model.Cart) (*model.Cart, error) {
	args := m.Called(cart)
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *CartRepositoryMock) DeleteCart(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
