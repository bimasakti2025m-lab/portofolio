package usecase

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type CartItemUseCaseMock struct {
	mock.Mock
}

func (m *CartItemUseCaseMock) CreateCartItem(cartItem model.CartItem) (*model.CartItem, error) {
	args := m.Called(cartItem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CartItem), args.Error(1)
}

func (m *CartItemUseCaseMock) GetAllCartItems() ([]model.CartItem, error) {
	args := m.Called()
	return args.Get(0).([]model.CartItem), args.Error(1)
}

func (m *CartItemUseCaseMock) GetCartItemByID(id int) (model.CartItem, error) {
	args := m.Called(id)
	return args.Get(0).(model.CartItem), args.Error(1)
}

func (m *CartItemUseCaseMock) UpdateCartItem(cartItem model.CartItem) (*model.CartItem, error) {
	args := m.Called(cartItem)
	return args.Get(0).(*model.CartItem), args.Error(1)
}

func (m *CartItemUseCaseMock) DeleteCartItem(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
