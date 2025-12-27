package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductUsecaseTestSuite struct {
	suite.Suite
	productRepo *repository.ProductRepositoryMock
	usecase ProductUsecase
}

func (suite *ProductUsecaseTestSuite) SetupTest() {
	suite.productRepo = new(repository.ProductRepositoryMock)
	suite.usecase = NewProductUsecase(suite.productRepo)
}
func TestProductUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUsecaseTestSuite))
}

func (suite *ProductUsecaseTestSuite) TestGetAllProducts_Success() {
	products := []model.Product{
		{ID: 1, Name: "Product 1", Description: "Description 1", Price: 100, Stock: 10},
		{ID: 2, Name: "Product 2", Description: "Description 2", Price: 200, Stock: 20},
	}
	suite.productRepo.On("GetAllProducts").Return(products, nil)
	result, err := suite.usecase.GetAllProducts()
	suite.NoError(err)
	suite.Equal(products, result)
}

func (suite *ProductUsecaseTestSuite) TestGetAllProducts_Failed() {
	suite.productRepo.On("GetAllProducts").Return([]model.Product(nil), assert.AnError)
	result, err := suite.usecase.GetAllProducts()
	suite.Error(err)
	suite.Nil(result)
}

func (suite *ProductUsecaseTestSuite) TestCreateProduct_Success() {
	product := model.Product{Name: "New Product", Description: "New Description", Price: 150, Stock: 15}
	createdProduct := product
	createdProduct.ID = 1
	suite.productRepo.On("CreateProduct", &product).Return(&createdProduct, nil)
	result, err := suite.usecase.CreateProduct(product)
	suite.NoError(err)
	suite.Equal(&createdProduct, result)
}

func (suite *ProductUsecaseTestSuite) TestCreateProduct_Failed() {
	product := model.Product{Name: "New Product", Description: "New Description", Price: 150, Stock: 15}
	suite.productRepo.On("CreateProduct", &product).Return((*model.Product)(nil), assert.AnError)
	result, err := suite.usecase.CreateProduct(product)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *ProductUsecaseTestSuite) TestDeleteProduct_Success() {
	suite.productRepo.On("DeleteProduct", 1).Return(nil)
	err := suite.usecase.DeleteProduct(1)
	suite.NoError(err)
}

func (suite *ProductUsecaseTestSuite) TestDeleteProduct_Failed() {
	suite.productRepo.On("DeleteProduct", 1).Return(assert.AnError)
	err := suite.usecase.DeleteProduct(1)
	suite.Error(err)
}

func (suite *ProductUsecaseTestSuite) TestGetProductByID_Success() {
	product := model.Product{ID: 1, Name: "Product 1", Description: "Description 1", Price: 100, Stock: 10}
	suite.productRepo.On("GetProductByID", 1).Return(product, nil)
	result, err := suite.usecase.GetProductByID(1)
	suite.NoError(err)
	suite.Equal(product, result)
}

func (suite *ProductUsecaseTestSuite) TestGetProductByID_Failed() {
	suite.productRepo.On("GetProductByID", 1).Return(model.Product{}, assert.AnError)
	result, err := suite.usecase.GetProductByID(1)
	suite.Error(err)
	suite.Equal(model.Product{}, result)
}

func (suite *ProductUsecaseTestSuite) TestUpdateProduct_Success() {
	product := &model.Product{ID: 1, Name: "Updated Product", Description: "Updated Description", Price: 120, Stock: 12}
	updatedProduct := *product
	suite.productRepo.On("UpdateProduct", product).Return(&updatedProduct, nil)
	result, err := suite.usecase.UpdateProduct(product)
	suite.NoError(err)
	suite.Equal(&updatedProduct, result)
}

func (suite *ProductUsecaseTestSuite) TestUpdateProduct_Failed() {
	product := &model.Product{ID: 1, Name: "Updated Product", Description: "Updated Description", Price: 120, Stock: 12}
	suite.productRepo.On("UpdateProduct", product).Return((*model.Product)(nil), assert.AnError)
	result, err := suite.usecase.UpdateProduct(product)
	suite.Error(err)
	suite.Nil(result)
}