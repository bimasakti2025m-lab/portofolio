package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type expensesRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	er      ExpenseRepository
}

var expectedExpense = entity.Expense{
	ID:              "uuid-expense-test",
	Date:            time.Now(),
	Amount:          10000,
	TransactionType: "CREDIT",
	Balance:         10000,
	Description:     "test",
	CreatedAt:       time.Now(),
	UpdatedAt:       time.Now(),
}

func TestExpensesRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(expensesRepositoryTestSuite))
}

func (s *expensesRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	s.NoError(err)

	s.mockDb = mockDb
	s.mockSql = mockSql
	s.er = NewExpenseRepository(mockDb)
}

func (s *expensesRepositoryTestSuite) TestCreate_success() {
	expenseRows := sqlmock.NewRows([]string{"id", "balance", "created_at"}).AddRow(expectedExpense.ID, expectedExpense.Balance, expectedExpense.CreatedAt)

	s.mockSql.ExpectQuery(regexp.QuoteMeta(`INSERT INTO expenses (date, amount, transaction_type, balance, description, user_id, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, balance, created_at`)).
		WithArgs(expectedExpense.Date, expectedExpense.Amount, expectedExpense.TransactionType, expectedExpense.Balance, expectedExpense.Description, expectedExpense.UserId, expectedExpense.UpdatedAt).
		WillReturnRows(
			expenseRows,
		)

	expense, err := s.er.Create(expectedExpense)

	s.Nil(err)
	s.Equal(expectedExpense.ID, expense.ID)
	s.Equal(expectedExpense.Balance, expense.Balance)
}

func (s *expensesRepositoryTestSuite) TestCreate_failed() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`INSERT INTO expenses (date, amount, transaction_type, balance, description, user_id, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, balance, created_at`)).
		WithArgs(expectedExpense.Date, expectedExpense.Amount, expectedExpense.TransactionType, expectedExpense.Balance, expectedExpense.Description, expectedExpense.UserId, expectedExpense.UpdatedAt).
		WillReturnError(sql.ErrConnDone)

	expense, err := s.er.Create(expectedExpense)

	s.NotNil(err)
	s.Equal(entity.Expense{}, expense)
}


func (s *expensesRepositoryTestSuite) TestGetBalance_success() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT balance FROM expenses WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`)).
		WithArgs("user-uuid-test").
		WillReturnRows(
			sqlmock.NewRows([]string{"balance"}).AddRow(15000),
		)

	balance, err := s.er.GetBalance("user-uuid-test")

	s.Nil(err)
	s.Equal(15000.0, balance)
}

func (s *expensesRepositoryTestSuite) TestGetBalance_failed() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT balance FROM expenses WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`)).
		WithArgs("user-uuid-test").
		WillReturnError(sql.ErrConnDone)

	balance, err := s.er.GetBalance("user-uuid-test")

	s.NotNil(err)
	s.Equal(0.0, balance)
}


func (s *expensesRepositoryTestSuite) TestGetBalance_noRows() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT balance FROM expenses WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`)).
		WithArgs("user-uuid-test").
		WillReturnError(sql.ErrNoRows)

	balance, err := s.er.GetBalance("user-uuid-test")

	s.NotNil(err)
	s.Equal(0.0, balance)
}


func (s *expensesRepositoryTestSuite) TestList_success() {
	// prepare
	expenseRows := sqlmock.NewRows([]string{"id", "date", "amount", "transaction_type", "balance", "description", "created_at", "updated_at"}).
		AddRow(
			expectedExpense.ID,
			expectedExpense.Date,
			expectedExpense.Amount,
			expectedExpense.TransactionType,
			expectedExpense.Balance,
			expectedExpense.Description,
			expectedExpense.CreatedAt,
			expectedExpense.UpdatedAt,
		)
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE user_id = $3 ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(10, 0, "user-uuid-test").
		WillReturnRows(expenseRows)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM expenses WHERE user_id=$1`)).
		WithArgs("user-uuid-test").
		WillReturnRows(countRows)

	// execute
	expenses, _, err := s.er.List(1, 10, "", "", "user-uuid-test")

	// assert
	s.Nil(err)
	s.Len(expenses, 1)
	s.Equal(expectedExpense.ID, expenses[0].ID)
}

func (s *expensesRepositoryTestSuite) TestList_failed(){
	// prepare
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE user_id = $3 ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(10, 0, "user-uuid-test").
		WillReturnError(sql.ErrConnDone)

	// execute
	expenses, _, err := s.er.List(1, 10, "", "", "user-uuid-test")

	// assert
	s.NotNil(err)
	s.Nil(expenses)
}

func (s *expensesRepositoryTestSuite) TestList_queryError() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE user_id = $3 ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(10, 0, "user-uuid-test").
		WillReturnError(sql.ErrConnDone)

	expenses, _, err := s.er.List(1, 10, "", "", "user-uuid-test")

	s.NotNil(err)
	s.Nil(expenses)
}

func (s *expensesRepositoryTestSuite) TestGetByTransaction_success() {
	expenseRows := sqlmock.NewRows([]string{"id", "date", "amount", "transaction_type", "balance", "description", "created_at", "updated_at"}).
		AddRow(
			expectedExpense.ID,
			expectedExpense.Date,
			expectedExpense.Amount,
			expectedExpense.TransactionType,
			expectedExpense.Balance,
			expectedExpense.Description,
			expectedExpense.CreatedAt,
			expectedExpense.UpdatedAt,
		)
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE transaction_type=$1 AND user_id = $2 ORDER BY created_at DESC`)).
		WithArgs("CREDIT", "user-uuid-test").
		WillReturnRows(expenseRows)

	expenses, err := s.er.GetByTransaction("CREDIT", "user-uuid-test")

	s.Nil(err)
	s.Len(expenses, 1)
	s.Equal(expectedExpense.ID, expenses[0].ID)
}

func (s *expensesRepositoryTestSuite) TestGetByTransaction_failed() {
	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE transaction_type=$1 AND user_id = $2 ORDER BY created_at DESC`)).
		WithArgs("CREDIT", "user-uuid-test").
		WillReturnError(sql.ErrConnDone)

	expenses, err := s.er.GetByTransaction("CREDIT", "user-uuid-test")

	s.NotNil(err)
	s.Nil(expenses)
}