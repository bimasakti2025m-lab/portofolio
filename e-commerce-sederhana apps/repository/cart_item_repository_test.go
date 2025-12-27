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

type cartItemRepositorySuite struct {
	suite.Suite
	cir     CartItemRepository
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
}

func TestCartItemRepositorySuite(t *testing.T) {
	suite.Run(t, new(cartItemRepositorySuite))
}

func (s *cartItemRepositorySuite) SetupTest() {
	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.mockDB = mockDB
	s.mockSQL = mockSQL
	s.cir = NewCartItemRepository(mockDB)
}

func (s *cartItemRepositorySuite) TestCreateCartItem_Success() {
	cartItem := model.CartItem{
		CartID:    1,
		ProductID: 1,
		Quantity:  2,
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs(cartItem.CartID, cartItem.ProductID, cartItem.Quantity).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := s.cir.CreateCartItem(&cartItem)
	s.NoError(err)
	s.Equal(1, result.ID)
}

func (s *cartItemRepositorySuite) TestCreateCartItem_Failed() {
	cartItem := model.CartItem{
		CartID:    1,
		ProductID: 1,
		Quantity:  2,
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs(cartItem.CartID, cartItem.ProductID, cartItem.Quantity).
		WillReturnError(errors.New("db error"))

	_, err := s.cir.CreateCartItem(&cartItem)
	s.Error(err)
}

func (s *cartItemRepositorySuite) TestGetAllCartItems_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, cart_id, product_id, quantity FROM cart_items")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity"}).
			AddRow(1, 1, 1, 2).
			AddRow(2, 1, 2, 1))

	cartItems, err := s.cir.GetAllCartItems()
	s.NoError(err)
	s.Len(cartItems, 2)
	s.Equal(1, cartItems[0].ProductID)
}

func (s *cartItemRepositorySuite) TestGetAllCartItems_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, cart_id, product_id, quantity FROM cart_items")).
		WillReturnError(errors.New("db error"))

	_, err := s.cir.GetAllCartItems()
	s.Error(err)
}

func (s *cartItemRepositorySuite) TestGetCartItemByID_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, cart_id, product_id, quantity FROM cart_items WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity"}).
			AddRow(1, 1, 1, 2))

	cartItem, err := s.cir.GetCartItemByID(1)
	s.NoError(err)
	s.Equal(1, cartItem.ID)
	s.Equal(2, cartItem.Quantity)
}

func (s *cartItemRepositorySuite) TestGetCartItemByID_NotFound() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, cart_id, product_id, quantity FROM cart_items WHERE id = $1")).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	cartItem, err := s.cir.GetCartItemByID(99)
	s.NoError(err) // Should return no error, but an empty cart item
	s.Equal(0, cartItem.ID)
}

func (s *cartItemRepositorySuite) TestGetCartItemByID_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, cart_id, product_id, quantity FROM cart_items WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err := s.cir.GetCartItemByID(1)
	s.Error(err)
}

func (s *cartItemRepositorySuite) TestUpdateCartItem_Success() {
	cartItem := model.CartItem{ID: 1, CartID: 1, ProductID: 1, Quantity: 3}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE cart_items SET cart_id = $2, product_id = $3, quantity = $4 WHERE id = $1")).
		WithArgs(cartItem.ID, cartItem.CartID, cartItem.ProductID, cartItem.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.cir.UpdateCartItem(&cartItem)
	s.NoError(err)
	s.Equal(cartItem.Quantity, result.Quantity)
}

func (s *cartItemRepositorySuite) TestUpdateCartItem_Failed() {
	cartItem := model.CartItem{ID: 1, CartID: 1, ProductID: 1, Quantity: 3}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE cart_items SET cart_id = $2, product_id = $3, quantity = $4 WHERE id = $1")).
		WithArgs(cartItem.ID, cartItem.CartID, cartItem.ProductID, cartItem.Quantity).
		WillReturnError(errors.New("db error"))

	_, err := s.cir.UpdateCartItem(&cartItem)
	s.Error(err)
}

func (s *cartItemRepositorySuite) TestDeleteCartItem_Success() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM cart_items WHERE id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.cir.DeleteCartItem(1)
	s.NoError(err)
}

func (s *cartItemRepositorySuite) TestDeleteCartItem_Failed() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM cart_items WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err := s.cir.DeleteCartItem(1)
	s.Error(err)
}
