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

type cartRepositorySuite struct {
	suite.Suite
	cr      CartRepository
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
}

func TestCartRepositorySuite(t *testing.T) {
	suite.Run(t, new(cartRepositorySuite))
}

func (s *cartRepositorySuite) SetupTest() {
	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.mockDB = mockDB
	s.mockSQL = mockSQL
	s.cr = NewCartRepository(mockDB)
}

func (s *cartRepositorySuite) TestCreateCart_Success() {
	cart := model.Cart{
		UserID: 1,
		Status: "active",
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO carts (user_id, status) VALUES ($1, $2) RETURNING id")).
		WithArgs(cart.UserID, cart.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := s.cr.CreateCart(&cart)
	s.NoError(err)
	s.Equal(1, result.ID)
}

func (s *cartRepositorySuite) TestCreateCart_Failed() {
	cart := model.Cart{
		UserID: 1,
		Status: "active",
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO carts (user_id, status) VALUES ($1, $2) RETURNING id")).
		WithArgs(cart.UserID, cart.Status).
		WillReturnError(errors.New("db error"))

	_, err := s.cr.CreateCart(&cart)
	s.Error(err)
}

func (s *cartRepositorySuite) TestGetAllCarts_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, status FROM carts")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status"}).
			AddRow(1, 1, "active").
			AddRow(2, 2, "completed"))

	carts, err := s.cr.GetAllCarts()
	s.NoError(err)
	s.Len(carts, 2)
	s.Equal("active", carts[0].Status)
}

func (s *cartRepositorySuite) TestGetAllCarts_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, status FROM carts")).
		WillReturnError(errors.New("db error"))

	_, err := s.cr.GetAllCarts()
	s.Error(err)
}

func (s *cartRepositorySuite) TestGetCartByID_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, status FROM carts WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status"}).
			AddRow(1, 1, "active"))

	cart, err := s.cr.GetCartByID(1)
	s.NoError(err)
	s.Equal(1, cart.ID)
	s.Equal("active", cart.Status)
}

func (s *cartRepositorySuite) TestGetCartByID_NotFound() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, status FROM carts WHERE id = $1")).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	cart, err := s.cr.GetCartByID(99)
	s.NoError(err) // Should return no error, but an empty cart
	s.Equal(0, cart.ID)
}

func (s *cartRepositorySuite) TestGetCartByID_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, status FROM carts WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err := s.cr.GetCartByID(1)
	s.Error(err)
}

func (s *cartRepositorySuite) TestUpdateCart_Success() {
	cart := model.Cart{ID: 1, UserID: 1, Status: "completed"}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE carts SET user_id = $2, status = $3 WHERE id = $1")).
		WithArgs(cart.ID, cart.UserID, cart.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.cr.UpdateCart(&cart)
	s.NoError(err)
	s.Equal(cart.Status, result.Status)
}

func (s *cartRepositorySuite) TestUpdateCart_Failed() {
	cart := model.Cart{ID: 1, UserID: 1, Status: "completed"}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE carts SET user_id = $2, status = $3 WHERE id = $1")).
		WithArgs(cart.ID, cart.UserID, cart.Status).
		WillReturnError(errors.New("db error"))

	_, err := s.cr.UpdateCart(&cart)
	s.Error(err)
}

func (s *cartRepositorySuite) TestDeleteCart_Success() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM carts WHERE id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.cr.DeleteCart(1)
	s.NoError(err)
}

func (s *cartRepositorySuite) TestDeleteCart_Failed() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM carts WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err := s.cr.DeleteCart(1)
	s.Error(err)
}
