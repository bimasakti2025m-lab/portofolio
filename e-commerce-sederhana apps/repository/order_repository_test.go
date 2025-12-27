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

type orderRepositorySuite struct {
	suite.Suite
	or      OrderRepository
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
}

func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(orderRepositorySuite))
}

func (s *orderRepositorySuite) SetupTest() {
	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.mockDB = mockDB
	s.mockSQL = mockSQL
	s.or = NewOrderRepository(mockDB)
}

func (s *orderRepositorySuite) TestCreateOrder_Success() {
	order := model.Order{
		UserID:                1,
		Total:                 150.00,
		StatusPesanan:         "pending",
		TransactionIDMidtrans: "midtrans-123",
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO orders (user_id, total, status_pesanan, transaction_id_midtrans) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs(order.UserID, order.Total, order.StatusPesanan, order.TransactionIDMidtrans).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := s.or.CreateOrder(order)
	s.NoError(err)
	s.Equal(1, result.ID)
}

func (s *orderRepositorySuite) TestCreateOrder_Failed() {
	order := model.Order{
		UserID:                1,
		Total:                 150.00,
		StatusPesanan:         "pending",
		TransactionIDMidtrans: "midtrans-123",
	}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO orders (user_id, total, status_pesanan, transaction_id_midtrans) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs(order.UserID, order.Total, order.StatusPesanan, order.TransactionIDMidtrans).
		WillReturnError(errors.New("db error"))

	_, err := s.or.CreateOrder(order)
	s.Error(err)
}

func (s *orderRepositorySuite) TestGetAllOrders_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, total, status_pesanan, transaction_id_midtrans FROM orders")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "total", "status_pesanan", "transaction_id_midtrans"}).
			AddRow(1, 1, 100.00, "pending", "midtrans-1").
			AddRow(2, 2, 200.00, "completed", "midtrans-2"))

	orders, err := s.or.GetAllOrders()
	s.NoError(err)
	s.Len(orders, 2)
	s.Equal("pending", orders[0].StatusPesanan)
}

func (s *orderRepositorySuite) TestGetAllOrders_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, total, status_pesanan, transaction_id_midtrans FROM orders")).
		WillReturnError(errors.New("db error"))

	_, err := s.or.GetAllOrders()
	s.Error(err)
}

func (s *orderRepositorySuite) TestGetOrderById_Success() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, total, status_pesanan, transaction_id_midtrans FROM orders WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "total", "status_pesanan", "transaction_id_midtrans"}).
			AddRow(1, 1, 100.00, "pending", "midtrans-1"))

	order, err := s.or.GetOrderById(1)
	s.NoError(err)
	s.Equal(1, order.ID)
	s.Equal("pending", order.StatusPesanan)
}

func (s *orderRepositorySuite) TestGetOrderById_NotFound() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, total, status_pesanan, transaction_id_midtrans FROM orders WHERE id = $1")).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	order, err := s.or.GetOrderById(99)
	s.NoError(err) // Should return no error, but an empty order
	s.Equal(0, order.ID)
}

func (s *orderRepositorySuite) TestGetOrderById_Failed() {
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, total, status_pesanan, transaction_id_midtrans FROM orders WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err := s.or.GetOrderById(1)
	s.Error(err)
}

func (s *orderRepositorySuite) TestUpdateOrder_Success() {
	order := model.Order{ID: 1, StatusPesanan: "paid"}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status_pesanan = $2 WHERE id = $1")).
		WithArgs(order.ID, order.StatusPesanan).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.or.UpdateOrder(order)
	s.NoError(err)
	s.Equal(order.StatusPesanan, result.StatusPesanan)
}

func (s *orderRepositorySuite) TestUpdateOrder_Failed() {
	order := model.Order{ID: 1, StatusPesanan: "paid"}
	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status_pesanan = $2 WHERE id = $1")).
		WithArgs(order.ID, order.StatusPesanan).
		WillReturnError(errors.New("db error"))

	_, err := s.or.UpdateOrder(order)
	s.Error(err)
}

func (s *orderRepositorySuite) TestDeleteOrder_Success() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM orders WHERE id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.or.DeleteOrder(1)
	s.NoError(err)
}

func (s *orderRepositorySuite) TestDeleteOrder_Failed() {
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM orders WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err := s.or.DeleteOrder(1)
	s.Error(err)
}
