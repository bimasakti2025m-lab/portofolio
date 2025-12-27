package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"enigmacamp.com/toko-enigma/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    ProductRepository
}

func (s *ProductRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)
	s.mockDB = db
	s.mockSQL = mock
	s.repo = NewProductRepository(s.mockDB)
}

func TestProductRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ProductRepoTestSuite))
}

func (s *ProductRepoTestSuite) TestCreate_Success() {
	payload := model.Product{Name: "Buku", Unit: "pcs", Stock: 10, Price: 5000}
	expectedID := 1

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_product (name, unit, stock,price) VALUES ($1,$2,$3,$4) RETURNING id")).
		WithArgs(payload.Name, payload.Unit, payload.Stock, payload.Price).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	actual, err := s.repo.Create(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedID, actual.Id)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestCreate_Fail() {
	payload := model.Product{Name: "Buku", Unit: "pcs", Stock: 10, Price: 5000}
	expectedError := errors.New("insert error")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_product (name, unit, stock,price) VALUES ($1,$2,$3,$4) RETURNING id")).
		WithArgs(payload.Name, payload.Unit, payload.Stock, payload.Price).
		WillReturnError(expectedError)

	actual, err := s.repo.Create(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), 0, actual.Id)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestFindAll_Success() {
	expectedProducts := []model.Product{
		{Id: 1, Name: "Buku", Unit: "pcs", Stock: 10, Price: 5000},
		{Id: 2, Name: "Pensil", Unit: "pcs", Stock: 20, Price: 1000},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "unit", "stock", "price"}).
		AddRow(expectedProducts[0].Id, expectedProducts[0].Name, expectedProducts[0].Unit, expectedProducts[0].Stock, expectedProducts[0].Price).
		AddRow(expectedProducts[1].Id, expectedProducts[1].Name, expectedProducts[1].Unit, expectedProducts[1].Stock, expectedProducts[1].Price)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT * FROM mst_product")).WillReturnRows(rows)

	actual, err := s.repo.FindAll()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedProducts, actual)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestFindAll_Fail() {
	expectedError := errors.New("query error")
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT * FROM mst_product")).WillReturnError(expectedError)

	actual, err := s.repo.FindAll()

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), actual)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestFindById_Success() {
	expectedProduct := model.Product{Id: 1, Name: "Buku", Unit: "pcs", Stock: 10, Price: 5000}

	rows := sqlmock.NewRows([]string{"id", "name", "unit", "stock", "price"}).
		AddRow(expectedProduct.Id, expectedProduct.Name, expectedProduct.Unit, expectedProduct.Stock, expectedProduct.Price)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT * FROM mst_product WHERE id = $1")).
		WithArgs(expectedProduct.Id).
		WillReturnRows(rows)

	actual, err := s.repo.FindById(expectedProduct.Id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedProduct, actual)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestFindById_Fail() {
	expectedError := sql.ErrNoRows

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT * FROM mst_product WHERE id = $1")).
		WithArgs(1).
		WillReturnError(expectedError)

	actual, err := s.repo.FindById(1)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Product{}, actual)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestUpdate_Success() {
	payload := model.Product{Id: 1, Name: "Buku Tulis", Unit: "pack", Stock: 5, Price: 25000}

	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE mst_product SET id=$2, name=$3, unit=$4, stock =$5, price=$6 WHERE id = $1")).
		WithArgs(payload.Id, payload.Id, payload.Name, payload.Unit, payload.Stock, payload.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	actual, err := s.repo.Update(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), payload, actual)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestUpdate_Fail() {
	payload := model.Product{Id: 1, Name: "Buku Tulis", Unit: "pack", Stock: 5, Price: 25000}
	expectedError := errors.New("update error")

	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE mst_product SET id=$2, name=$3, unit=$4, stock =$5, price=$6 WHERE id = $1")).
		WithArgs(payload.Id, payload.Id, payload.Name, payload.Unit, payload.Stock, payload.Price).
		WillReturnError(expectedError)

	actual, err := s.repo.Update(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Product{}, actual)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestDelete_Success() {
	productID := 1

	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM mst_product WHERE id =$1")).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Delete(productID)

	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *ProductRepoTestSuite) TestDelete_Fail() {
	productID := 1
	expectedError := errors.New("delete error")

	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM mst_product WHERE id =$1")).
		WithArgs(productID).
		WillReturnError(expectedError)

	err := s.repo.Delete(productID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}
