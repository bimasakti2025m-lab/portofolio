package usecase

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type CartUseCaseMock struct {
	mock.Mock
}

func (m *CartUseCaseMock) CreateCart(cart model.Cart) (*model.Cart, error) {
	args := m.Called(cart)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *CartUseCaseMock) GetAllCarts() ([]model.Cart, error) {
	args := m.Called()
	return args.Get(0).([]model.Cart), args.Error(1)
}

func (m *CartUseCaseMock) GetCartByID(id int) (model.Cart, error) {
	args := m.Called(id)
	return args.Get(0).(model.Cart), args.Error(1)
}

func (m *CartUseCaseMock) UpdateCart(cart model.Cart) (*model.Cart, error) {
	args := m.Called(cart)
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *CartUseCaseMock) DeleteCart(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
