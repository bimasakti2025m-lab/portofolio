package usecase

import (
	"enigmacamp.com/toko-enigma/model"
	"github.com/stretchr/testify/mock"
)

type CartUseCaseMock struct {
	mock.Mock
}

func (m *CartUseCaseMock) CreateNewCart(cart model.Cart) (model.Cart, error) {
	args := m.Called(cart)
	return args.Get(0).(model.Cart), args.Error(1)
}

func (m *CartUseCaseMock) GetAllCart() ([]model.Cart, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Cart), args.Error(1)
}

func (m *CartUseCaseMock) GetCartById(id int) (model.Cart, error) {
	args := m.Called(id)
	return args.Get(0).(model.Cart), args.Error(1)
}

func (m *CartUseCaseMock) UpdateCartById(cart model.Cart) (model.Cart, error) {
	args := m.Called(cart)
	return args.Get(0).(model.Cart), args.Error(1)
}

func (m *CartUseCaseMock) DeleteCartById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
