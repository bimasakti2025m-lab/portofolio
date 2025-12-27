package usecase

import (
	"enigmacamp.com/toko-enigma/model"
	"github.com/stretchr/testify/mock"
)

type ProductUseCaseMock struct {
	mock.Mock
}

func (m *ProductUseCaseMock) CreateNewProduct(product model.Product) (model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *ProductUseCaseMock) GetAllProduct() ([]model.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *ProductUseCaseMock) GetProductById(id int) (model.Product, error) {
	args := m.Called(id)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *ProductUseCaseMock) UpdateProductById(product model.Product) (model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *ProductUseCaseMock) DeleteProductById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
