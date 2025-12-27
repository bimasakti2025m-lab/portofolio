package repository_test

import (
	"E-commerce-Sederhana/model"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	. "E-commerce-Sederhana/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type productRepositorySuite struct {
	suite.Suite
	pr      ProductRepository
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
}

func TestProductRepositorySuite(t *testing.T) {
	suite.Run(t, new(productRepositorySuite))
}

func (s *productRepositorySuite) SetupTest() {
	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.mockDB = mockDB
	s.mockSQL = mockSQL
	s.pr = NewProductRepository(mockDB)
}

func (s *productRepositorySuite) TestCreateProduct_Success() {
	product := model.Product{
		Name:        "Test Product",
		Description: "Description for test product",
		Price:       100.00,
		Stock:       10,
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs(product.Name, product.Description, product.Price, product.Stock).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := s.pr.CreateProduct(&product)
	s.NoError(err)
	s.Equal(1, result.ID)
}

func (s *productRepositorySuite) TestCreateProduct_Failed() {
	product := model.Product{
		Name:        "Test Product",
		Description: "Description for test product",
		Price:       100.00,
		Stock:       10,
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs(product.Name, product.Description, product.Price, product.Stock).
		WillReturnError(errors.New("db error"))

	_, err := s.pr.CreateProduct(&product)
	s.Error(err)
}

func (s *productRepositorySuite) TestGetAllProducts_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, price, stock FROM products")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "stock"}).
			AddRow(1, "Product 1", "Desc 1", 10.00, 5).
			AddRow(2, "Product 2", "Desc 2", 20.00, 15))

	products, err := s.pr.GetAllProducts()
	s.NoError(err)
	s.Len(products, 2)
	s.Equal("Product 1", products[0].Name)
}

func (s *productRepositorySuite) TestGetAllProducts_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, price, stock FROM products")).
		WillReturnError(errors.New("db error"))

	_, err := s.pr.GetAllProducts()
	s.Error(err)
}

func (s *productRepositorySuite) TestGetProductByID_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, price, stock FROM products WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "stock"}).
			AddRow(1, "Product 1", "Desc 1", 10.00, 5))

	product, err := s.pr.GetProductByID(1)
	s.NoError(err)
	s.Equal(1, product.ID)
	s.Equal("Product 1", product.Name)
}

func (s *productRepositorySuite) TestGetProductByID_NotFound() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, price, stock FROM products WHERE id = $1")).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	product, err := s.pr.GetProductByID(99)
	s.NoError(err) // Should return no error, but an empty product
	s.Equal(0, product.ID)
}

func (s *productRepositorySuite) TestGetProductByID_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, price, stock FROM products WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err := s.pr.GetProductByID(1)
	s.Error(err)
}

func (s *productRepositorySuite) TestUpdateProduct_Success() {
	product := model.Product{ID: 1, Name: "Updated Product", Description: "Updated Desc", Price: 150.00, Stock: 12}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE products SET name = $2, description = $3, price = $4, stock = $5 WHERE id = $1")).
		WithArgs(product.ID, product.Name, product.Description, product.Price, product.Stock).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.pr.UpdateProduct(&product)
	s.NoError(err)
	s.Equal(product.Name, result.Name)
}

func (s *productRepositorySuite) TestUpdateProduct_Failed() {
	product := model.Product{ID: 1, Name: "Updated Product", Description: "Updated Desc", Price: 150.00, Stock: 12}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE products SET name = $2, description = $3, price = $4, stock = $5 WHERE id = $1")).
		WithArgs(product.ID, product.Name, product.Description, product.Price, product.Stock).
		WillReturnError(errors.New("db error"))

	_, err := s.pr.UpdateProduct(&product)
	s.Error(err)
}

func (s *productRepositorySuite) TestDeleteProduct_Success() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.pr.DeleteProduct(1)
	s.NoError(err)
}

func (s *productRepositorySuite) TestDeleteProduct_Failed() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err := s.pr.DeleteProduct(1)
	s.Error(err)
}
