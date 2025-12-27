package repository

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type CartItemRepositoryMock struct {
	mock.Mock
}

func (m *CartItemRepositoryMock) CreateCartItem(cartItem *model.CartItem) (*model.CartItem, error) {
	args := m.Called(cartItem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CartItem), args.Error(1)
}

func (m *CartItemRepositoryMock) GetAllCartItems() ([]model.CartItem, error) {
	args := m.Called()
	return args.Get(0).([]model.CartItem), args.Error(1)
}

func (m *CartItemRepositoryMock) GetCartItemByID(id int) (model.CartItem, error) {
	args := m.Called(id)
	return args.Get(0).(model.CartItem), args.Error(1)
}

func (m *CartItemRepositoryMock) UpdateCartItem(cartItem *model.CartItem) (*model.CartItem, error) {
	args := m.Called(cartItem)
	return args.Get(0).(*model.CartItem), args.Error(1)
}

func (m *CartItemRepositoryMock) DeleteCartItem(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
