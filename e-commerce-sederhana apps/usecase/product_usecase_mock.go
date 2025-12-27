package usecase

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type ProductUseCaseMock struct {
	mock.Mock
}

func (m *ProductUseCaseMock) CreateProduct(product model.Product) (*model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *ProductUseCaseMock) GetAllProducts() ([]model.Product, error) {
	args := m.Called()
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *ProductUseCaseMock) GetProductByID(id int) (model.Product, error) {
	args := m.Called(id)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *ProductUseCaseMock) UpdateProduct(product *model.Product) (*model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *ProductUseCaseMock) DeleteProduct(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
