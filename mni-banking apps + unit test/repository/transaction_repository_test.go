package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"enigmacamp.com/mini-banking/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    TransactionRepository
}

func (s *TransactionRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)
	s.mockDB = db
	s.mockSQL = mock
	s.repo = NewTransactionRepository(s.mockDB)
}

func TestTransactionRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepoTestSuite))
}

func (s *TransactionRepoTestSuite) TestCreate_Success() {
	payload := model.Transaction{FromUserID: 1, ToUserID: 2, Amount: 100, Type: "transfer", Status: "completed"}
	expectedID := uint(1)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_transaction (from_user_id, to_user_id, amount, type, status) VALUES ($1, $2, $3, $4, $5) RETURNING id")).
		WithArgs(payload.FromUserID, payload.ToUserID, payload.Amount, payload.Type, payload.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	transaction, err := s.repo.Create(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedID, transaction.ID)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestCreate_Fail() {
	payload := model.Transaction{FromUserID: 1, ToUserID: 2, Amount: 100, Type: "transfer", Status: "completed"}
	expectedError := errors.New("failed to save transaction")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_transaction (from_user_id, to_user_id, amount, type, status) VALUES ($1, $2, $3, $4, $5) RETURNING id")).
		WithArgs(payload.FromUserID, payload.ToUserID, payload.Amount, payload.Type, payload.Status).
		WillReturnError(expectedError)

	transaction, err := s.repo.Create(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Transaction{}, transaction)
	assert.Equal(s.T(), expectedError, err)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestList_Success() {
	expectedTransaction := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100, Type: "transfer"}

	rows := sqlmock.NewRows([]string{"id", "from_user_id", "to_user_id", "amount", "type"}).
		AddRow(expectedTransaction.ID, expectedTransaction.FromUserID, expectedTransaction.ToUserID, expectedTransaction.Amount, expectedTransaction.Type)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction")).
		WillReturnRows(rows)

	transactions, err := s.repo.List()

	assert.NoError(s.T(), err)
	assert.Len(s.T(), transactions, 1)
	assert.Equal(s.T(), expectedTransaction, transactions[0])
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestList_Fail() {
	expectedError := errors.New("failed to retrieve list transaction")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction")).
		WillReturnError(expectedError)

	transactions, err := s.repo.List()

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), transactions)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestGet_Success() {
	expectedTransaction := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100, Type: "transfer"}

	rows := sqlmock.NewRows([]string{"id", "from_user_id", "to_user_id", "amount", "type"}).
		AddRow(expectedTransaction.ID, expectedTransaction.FromUserID, expectedTransaction.ToUserID, expectedTransaction.Amount, expectedTransaction.Type)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction WHERE id = $1")).
		WithArgs(uint32(1)).
		WillReturnRows(rows)

	transaction, err := s.repo.Get(uint32(1))

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedTransaction, transaction)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestGet_Fail() {
	expectedError := errors.New("failed to get transaction by ID")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction WHERE id = $1")).
		WithArgs(uint32(1)).
		WillReturnError(expectedError)

	transaction, err := s.repo.Get(uint32(1))

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Transaction{}, transaction)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestGetByUserId_Success() {
	userId := uint32(1)
	expectedTransaction := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100, Type: "transfer"}

	rows := sqlmock.NewRows([]string{"id", "from_user_id", "to_user_id", "amount", "type"}).
		AddRow(expectedTransaction.ID, expectedTransaction.FromUserID, expectedTransaction.ToUserID, expectedTransaction.Amount, expectedTransaction.Type)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction WHERE from_user_id = $1 OR to_user_id = $1")).
		WithArgs(userId).
		WillReturnRows(rows)

	transactions, err := s.repo.GetByUserId(userId)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), transactions, 1)
	assert.Equal(s.T(), expectedTransaction, transactions[0])
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestGetByUserId_Fail() {
	userId := uint32(1)
	expectedError := errors.New("failed to get transaction by user ID")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction WHERE from_user_id = $1 OR to_user_id = $1")).
		WithArgs(userId).
		WillReturnError(expectedError)

	transactions, err := s.repo.GetByUserId(userId)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), transactions)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestUpdate_Success() {
	payload := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 150, Type: "transfer"}

	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE mst_transaction SET from_user_id = $1, to_user_id = $2, amount = $3, type = $4 WHERE id = $5")).
		WithArgs(payload.FromUserID, payload.ToUserID, payload.Amount, payload.Type, payload.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	transaction, err := s.repo.Update(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), payload, transaction)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestUpdate_Fail() {
	payload := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 150, Type: "transfer"}
	expectedError := errors.New("failed to update transaction")

	s.mockSQL.ExpectExec(regexp.QuoteMeta("UPDATE mst_transaction SET from_user_id = $1, to_user_id = $2, amount = $3, type = $4 WHERE id = $5")).
		WithArgs(payload.FromUserID, payload.ToUserID, payload.Amount, payload.Type, payload.ID).
		WillReturnError(expectedError)

	transaction, err := s.repo.Update(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Transaction{}, transaction)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestDelete_Success() {
	id := uint32(1)
	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM mst_transaction WHERE id = $1")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Delete(id)

	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *TransactionRepoTestSuite) TestDelete_Fail() {
	id := uint32(1)
	expectedError := errors.New("failed to delete transaction")

	s.mockSQL.ExpectExec(regexp.QuoteMeta("DELETE FROM mst_transaction WHERE id = $1")).
		WithArgs(id).
		WillReturnError(expectedError)

	err := s.repo.Delete(id)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}
