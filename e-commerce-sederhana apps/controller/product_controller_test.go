package controller

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/usecase"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductControllerTestSuite struct {
	suite.Suite
	productUseCase *usecase.ProductUseCaseMock
	controller     *ProductController
	router         *gin.Engine
}

func (suite *ProductControllerTestSuite) SetupTest() {
	suite.productUseCase = new(usecase.ProductUseCaseMock)
	suite.controller = &ProductController{
		productUc: suite.productUseCase,
	}

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	suite.router.POST("/products", suite.controller.createProductHandler)
	suite.router.GET("/products", suite.controller.getAllProductHandler)
	suite.router.GET("/products/:id", suite.controller.getProductByIdHandler)
	suite.router.PUT("/products/:id", suite.controller.updateProductHandler)
	suite.router.DELETE("/products/:id", suite.controller.deleteProductHandler)
}

func TestProductControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductControllerTestSuite))
}

func (suite *ProductControllerTestSuite) TestCreateProduct_Success() {
	product := model.Product{Name: "Laptop", Price: 10000000}
	createdProduct := product
	createdProduct.ID = 1

	suite.productUseCase.On("CreateProduct", product).Return(&createdProduct, nil)

	body, _ := json.Marshal(product)
	req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response model.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), createdProduct.ID, response.ID)
}

func (suite *ProductControllerTestSuite) TestGetAllProducts_Success() {
	products := []model.Product{
		{ID: 1, Name: "Laptop", Price: 10000000},
		{ID: 2, Name: "Mouse", Price: 100000},
	}

	suite.productUseCase.On("GetAllProducts").Return(products, nil)

	req, _ := http.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response []model.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(suite.T(), response, 2)
}

func (suite *ProductControllerTestSuite) TestGetProductById_Success() {
	product := model.Product{ID: 1, Name: "Laptop", Price: 10000000}

	suite.productUseCase.On("GetProductByID", 1).Return(product, nil)

	req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response model.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), product.Name, response.Name)
}

func (suite *ProductControllerTestSuite) TestGetProductById_NotFound() {
	suite.productUseCase.On("GetProductByID", 99).Return(model.Product{}, errors.New("product not found"))

	req, _ := http.NewRequest(http.MethodGet, "/products/99", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *ProductControllerTestSuite) TestUpdateProduct_Success() {
	product := model.Product{Name: "Laptop Updated", Price: 12000000}
	updatedProduct := product
	updatedProduct.ID = 1

	// Matcher untuk pointer product
	suite.productUseCase.On("UpdateProduct", mock.MatchedBy(func(p *model.Product) bool {
		return p.Name == product.Name
	})).Return(&updatedProduct, nil)

	body, _ := json.Marshal(product)
	req, _ := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response model.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Laptop Updated", response.Name)
}

func (suite *ProductControllerTestSuite) TestDeleteProduct_Success() {
	suite.productUseCase.On("DeleteProduct", 1).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *ProductControllerTestSuite) TestDeleteProduct_Error() {
	suite.productUseCase.On("DeleteProduct", 1).Return(errors.New("delete failed"))

	req, _ := http.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
