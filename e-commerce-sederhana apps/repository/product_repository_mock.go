package repository

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) GetAllProducts() ([]model.Product, error) {
	args := m.Called()
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetProductByID(id int) (model.Product, error) {
	args := m.Called(id)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *ProductRepositoryMock) CreateProduct(product *model.Product) (*model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *ProductRepositoryMock) UpdateProduct(product *model.Product) (*model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *ProductRepositoryMock) DeleteProduct(id int) error {
	args := m.Called(id)
	return args.Error(0)

}
