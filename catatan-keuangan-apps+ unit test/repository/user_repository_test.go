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

var expectedUser = entity.User{
	ID:        "uuid-test",
	Username:  "test",
	Role:      "user",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Expenses: []entity.Expense{
		{
			ID:              "uuid-expense-test",
			Date:            time.Now(),
			Amount:          10000,
			TransactionType: "CREDIT",
			Balance:         10000,
			Description:     "test",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	},
}

type userRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	ur      UserRepository
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(userRepositoryTestSuite))
}

func (s *userRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	s.NoError(err)

	s.mockDb = mockDb
	s.mockSql = mockSql
	s.ur = NewUserRepository(mockDb)
}

func (s *userRepositoryTestSuite) TestGet_success() {
	userRows := sqlmock.NewRows([]string{"id", "username", "role", "created_at", "updated_at"}).AddRow(
		expectedUser.ID,
		expectedUser.Username,
		expectedUser.Role,
		expectedUser.CreatedAt,
		expectedUser.UpdatedAt,
	)

	s.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role, created_at, updated_at FROM users WHERE id=$1")).
		WithArgs(expectedUser.ID).WillReturnRows(
		userRows,
	)

	expenseRows := sqlmock.NewRows([]string{"user_id", "id", "date", "amount", "transaction_type", "balance", "description", "created_at", "updated_at"}).AddRow(
		expectedUser.ID,
		expectedUser.Expenses[0].ID,
		expectedUser.Expenses[0].Date,
		expectedUser.Expenses[0].Amount,
		expectedUser.Expenses[0].TransactionType,
		expectedUser.Expenses[0].Balance,
		expectedUser.Expenses[0].Description,
		expectedUser.Expenses[0].CreatedAt,
		expectedUser.Expenses[0].UpdatedAt,
	)

	s.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, e.id, e.date, e.amount, e.transaction_type,e.balance, e.description, e.created_at, e.updated_at
FROM users u JOIN expenses e on u.id = e.user_id WHERE u.id = $1`)).WithArgs(expectedUser.ID).WillReturnRows(
		expenseRows,
	)

	user, err := s.ur.Get("uuid-test")

	s.Nil(err)
	s.Equal(expectedUser, user)
}
