package repository

import (
	"database/sql"
	"errors"

	"regexp"
	"testing"

	"enigmacamp.com/url-shortener/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo UserRepository
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
	expectedID := uint32(1)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_user (username, password, role) VALUES  ($1, $2, $3) RETURNING id")).
		WithArgs(payload.Username, payload.Password, payload.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	user, err := s.repo.Create(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedID, user.Id)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestCreate_Fail() {
	payload := model.UserCredential{Username: "test", Password: "password", Role: "user"}
	expectedError := errors.New("failed to save user")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_user (username, password, role) VALUES  ($1, $2, $3) RETURNING id")).
		WithArgs(payload.Username, payload.Password, payload.Role).
		WillReturnError(expectedError)

	user, err := s.repo.Create(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), uint32(0), user.Id)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}



func (s *UserRepoTestSuite) TestList_Success() {
	expectedUsers := []model.UserCredential{
		{Id: 1, Username: "admin", Role: "admin"},
		{Id: 2, Username: "user1", Role: "user"},
	}

	rows := sqlmock.NewRows([]string{"id", "username", "role"}).
		AddRow(expectedUsers[0].Id, expectedUsers[0].Username, expectedUsers[0].Role).
		AddRow(expectedUsers[1].Id, expectedUsers[1].Username, expectedUsers[1].Role)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role FROM mst_user")).
		WillReturnRows(rows)

	users, err := s.repo.List()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUsers, users)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestList_Fail() {
	expectedError := errors.New("failed to retrieve list user")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role FROM mst_user")).
		WillReturnError(expectedError)

	users, err := s.repo.List()

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), users)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestGet_Success() {
	expectedUser := model.UserCredential{Id: 1, Username: "admin", Role: "admin"}

	rows := sqlmock.NewRows([]string{"id", "username", "role"}).
		AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Role)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role FROM mst_user WHERE id = $1")).
		WithArgs(expectedUser.Id).
		WillReturnRows(rows)

	user, err := s.repo.Get(expectedUser.Id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, user)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestGet_Fail() {
	expectedError := errors.New("failed to get user by ID")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role FROM mst_user WHERE id = $1")).
		WithArgs(uint32(1)).
		WillReturnError(expectedError)

	user, err := s.repo.Get(uint32(1))

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.UserCredential{}, user)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestGetByUsernamePassword_Success() {
	expectedUser := model.UserCredential{Id: 1, Username: "admin", Password: "admin", Role: "admin"}

	rows := sqlmock.NewRows([]string{"id", "username", "role"}).
		AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Role)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role FROM mst_user WHERE username = $1 and password = $2")).
		WithArgs(expectedUser.Username, expectedUser.Password).
		WillReturnRows(rows)

	user, err := s.repo.GetByUsernamePassword(expectedUser.Username, expectedUser.Password)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser.Id, user.Id)
	assert.Equal(s.T(), expectedUser.Username, user.Username)
	assert.Equal(s.T(), expectedUser.Role, user.Role)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UserRepoTestSuite) TestGetByUsernamePassword_Fail() {
	expectedError := errors.New("failed to get user by username and password")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role FROM mst_user WHERE username = $1 and password = $2")).
		WithArgs("admin", "admin").
		WillReturnError(expectedError)

	user, err := s.repo.GetByUsernamePassword("admin", "admin")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.UserCredential{}, user)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}


