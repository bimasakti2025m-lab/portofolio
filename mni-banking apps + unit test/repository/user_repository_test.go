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

type UserRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    UserRepository
}

func (s *UserRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)
	s.mockDB = db
	s.mockSQL = mock
	s.repo = NewUserRepository(s.mockDB)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}

func (s *UserRepoTestSuite) TestCreate_Success() {
	payload := model.UserCredential{Username: "test", Password: "password", Role: "user"}
	expectedUser := model.UserCredential{Id: 1, Username: "test", Password: "password", Role: "user", Balance: 0}

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_user (username, password, role, balance) VALUES  ($1, $2, $3, $4) RETURNING id")).
		WithArgs(payload.Username, payload.Password, payload.Role, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	user, err := s.repo.Create(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, user)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestCreate_Fail() {
	payload := model.UserCredential{Username: "test", Password: "password", Role: "user"}
	expectedError := errors.New("failed to save user")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_user (username, password, role, balance) VALUES  ($1, $2, $3, $4) RETURNING id")).
		WithArgs(payload.Username, payload.Password, payload.Role, 0).
		WillReturnError(expectedError)

	user, err := s.repo.Create(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.UserCredential{}, user)
	assert.Equal(s.T(), expectedError, err)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestList_Success() {
	expectedUser := model.UserCredential{Id: 1, Username: "admin", Role: "admin", Balance: 1000}

	rows := sqlmock.NewRows([]string{"id", "username", "role", "balance"}).
		AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Role, expectedUser.Balance)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role, balance FROM mst_user")).
		WillReturnRows(rows)

	users, err := s.repo.List()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, users[0])
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestList_Fail() {
	expectedError := errors.New("failed to retrieve list user")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role, balance FROM mst_user")).
		WillReturnError(expectedError)

	users, err := s.repo.List()

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), users)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestGet_Success() {
	expectedUser := model.UserCredential{Id: 1, Username: "admin", Role: "admin", Balance: 1000}

	rows := sqlmock.NewRows([]string{"id", "username", "role", "balance"}).
		AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Role, expectedUser.Balance)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role, balance FROM mst_user WHERE id = $1")).
		WithArgs(uint32(1)).
		WillReturnRows(rows)

	user, err := s.repo.Get(uint32(1))

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser.Id, user.Id)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestGet_Fail() {
	expectedError := errors.New("failed to get user by ID")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role, balance FROM mst_user WHERE id = $1")).
		WithArgs(uint32(1)).
		WillReturnError(expectedError)

	user, err := s.repo.Get(uint32(1))

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.UserCredential{}, user)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestGetByUsername_success() {
	expectedUser := model.UserCredential{Id: 1, Username: "admin", Password: "password", Role: "admin", Balance: 1000}

	rows := sqlmock.NewRows([]string{"id", "username", "password", "role", "balance"}).
		AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Password, expectedUser.Role, expectedUser.Balance)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, password, role, balance FROM mst_user WHERE username = $1")).
		WithArgs(expectedUser.Username).
		WillReturnRows(rows)

	user, err := s.repo.GetByUsername(expectedUser.Username)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, user)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestGetByUsername_fail() {
	expectedError := errors.New("failed to get user by username")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, password, role, balance FROM mst_user WHERE username = $1")).
		WithArgs("admin").
		WillReturnError(expectedError)

	user, err := s.repo.GetByUsername("admin")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.UserCredential{}, user)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}


