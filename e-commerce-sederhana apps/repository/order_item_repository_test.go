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

type orderItemRepositorySuite struct {
	suite.Suite
	oir     OrderItemRepository
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
}

func TestOrderItemRepositorySuite(t *testing.T) {
	suite.Run(t, new(orderItemRepositorySuite))
}

func (s *orderItemRepositorySuite) SetupTest() {
	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.mockDB = mockDB
	s.mockSQL = mockSQL
	s.oir = NewOrderItemRepository(mockDB)
}

func (s *orderItemRepositorySuite) TestCreateOrderItem_Success() {
	orderItem := model.OrderItem{
		OrderID:       1,
		ProductID:     1,
		Quantity:      2,
		PriceSnapshot: 50.00,
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO order_items (order_id, product_id, quantity, price_snapshot) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs(orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.PriceSnapshot).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := s.oir.CreateOrderItem(&orderItem)
	s.NoError(err)
	s.Equal(1, result.ID)
}

func (s *orderItemRepositorySuite) TestCreateOrderItem_Failed() {
	orderItem := model.OrderItem{
		OrderID:       1,
		ProductID:     1,
		Quantity:      2,
		PriceSnapshot: 50.00,
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO order_items (order_id, product_id, quantity, price_snapshot) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs(orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.PriceSnapshot).
		WillReturnError(errors.New("db error"))

	_, err := s.oir.CreateOrderItem(&orderItem)
	s.Error(err)
}

func (s *orderItemRepositorySuite) TestGetAllOrderItems_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, product_id, quantity, price_snapshot FROM order_items")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price_snapshot"}).
			AddRow(1, 1, 1, 2, 50.00).
			AddRow(2, 1, 2, 1, 100.00))

	orderItems, err := s.oir.GetAllOrderItems()
	s.NoError(err)
	s.Len(orderItems, 2)
	s.Equal(1, orderItems[0].ProductID)
}

func (s *orderItemRepositorySuite) TestGetAllOrderItems_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, product_id, quantity, price_snapshot FROM order_items")).
		WillReturnError(errors.New("db error"))

	_, err := s.oir.GetAllOrderItems()
	s.Error(err)
}

func (s *orderItemRepositorySuite) TestGetOrderItemByID_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, product_id, quantity, price_snapshot FROM order_items WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price_snapshot"}).
			AddRow(1, 1, 1, 2, 50.00))

	orderItem, err := s.oir.GetOrderItemByID(1)
	s.NoError(err)
	s.Equal(1, orderItem.ID)
	s.Equal(2, orderItem.Quantity)
}

func (s *orderItemRepositorySuite) TestGetOrderItemByID_NotFound() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, product_id, quantity, price_snapshot FROM order_items WHERE id = $1")).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	orderItem, err := s.oir.GetOrderItemByID(99)
	s.NoError(err) // Should return no error, but an empty order item
	s.Equal(0, orderItem.ID)
}

func (s *orderItemRepositorySuite) TestGetOrderItemByID_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, product_id, quantity, price_snapshot FROM order_items WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err := s.oir.GetOrderItemByID(1)
	s.Error(err)
}

func (s *orderItemRepositorySuite) TestUpdateOrderItem_Success() {
	orderItem := model.OrderItem{ID: 1, OrderID: 1, ProductID: 1, Quantity: 3, PriceSnapshot: 55.00}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE order_items SET order_id = $2, product_id = $3, quantity = $4, price_snapshot = $5 WHERE id = $1")).
		WithArgs(orderItem.ID, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.PriceSnapshot).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.oir.UpdateOrderItem(&orderItem)
	s.NoError(err)
	s.Equal(orderItem.Quantity, result.Quantity)
}

func (s *orderItemRepositorySuite) TestUpdateOrderItem_Failed() {
	orderItem := model.OrderItem{ID: 1, OrderID: 1, ProductID: 1, Quantity: 3, PriceSnapshot: 55.00}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE order_items SET order_id = $2, product_id = $3, quantity = $4, price_snapshot = $5 WHERE id = $1")).
		WithArgs(orderItem.ID, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.PriceSnapshot).
		WillReturnError(errors.New("db error"))

	_, err := s.oir.UpdateOrderItem(&orderItem)
	s.Error(err)
}

func (s *orderItemRepositorySuite) TestDeleteOrderItem_Success() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM order_items WHERE id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.oir.DeleteOrderItem(1)
	s.NoError(err)
}

func (s *orderItemRepositorySuite) TestDeleteOrderItem_Failed() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM order_items WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err := s.oir.DeleteOrderItem(1)
	s.Error(err)
}
