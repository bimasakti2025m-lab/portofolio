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

type userRepositorySuite struct {
	suite.Suite
	u       UserRepository
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(userRepositorySuite))
} 

// membuat method setup
func (u *userRepositorySuite) SetupTest() {

	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		u.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	u.mockDB = mockDB
	u.mockSQL = mockSQL
	u.u = NewUserRepository(mockDB)
}

func (u *userRepositorySuite) TestCreate_Success() {
	user := model.User{
		Username: "username",
		Email:    "email",
		Password: "password",
		Role:     "user",
	}

	u.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs(user.Username, user.Email, user.Password, user.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := u.u.Create(&user)
	u.NoError(err)
	u.Equal(1, result.ID)
}


func (u *userRepositorySuite) TestCreate_UsernameAlreadyExists() {
	user := model.User{
		Username: "username",
		Password: "password",
		Role:     "user",
	}

	u.mockSQL.ExpectQuery("INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id").
		WithArgs(user.Username, user.Password, user.Role).
		WillReturnError(errors.New("username already exists"))

	_, err := u.u.Create(&user)
	u.Error(err)
}

func (u *userRepositorySuite) TestCreate_Failed() {
	user := model.User{
		Username: "username",
		Password: "password",
		Role:     "user",
	}

	u.mockSQL.ExpectQuery("INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id").
		WithArgs(user.Username, user.Password, user.Role).
		WillReturnError(errors.New("error"))

	_, err := u.u.Create(&user)
	u.Error(err)
}

func (u *userRepositorySuite) TestGetAllUsers_Success() {

	u.mockSQL.ExpectQuery("SELECT id, username, email, password, role FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "role"}).
			AddRow(1, "username", "email", "password", "user"))

	users, err := u.u.GetAllUsers()
	u.NoError(err)
	u.Equal(1, len(users))
}

func (u *userRepositorySuite) TestGetAllUsers_Failed() {

	u.mockSQL.ExpectQuery("SELECT id, username, password, role FROM users").
		WillReturnError(errors.New("error"))

	_, err := u.u.GetAllUsers()
	u.Error(err)
}

func (u *userRepositorySuite) TestGetUserByUsername_Success() {

	u.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, username, email, password, role FROM users WHERE username = $1")).
		WithArgs("username").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "role"}).
			AddRow(1, "username", "email", "password", "user"))

	user, err := u.u.GetUserByUsername("username")
	u.NoError(err)
	u.Equal(1, user.ID)
}

func (u *userRepositorySuite) TestGetUserByUsername_UserNotFound() {

	u.mockSQL.ExpectQuery("SELECT id, username, email, password, role FROM users WHERE username = $1").
		WithArgs("username").
		WillReturnError(sql.ErrNoRows)

	_, err := u.u.GetUserByUsername("username")
	u.Error(err)
}

func (u *userRepositorySuite) TestGetUserByUsername_Failed() {

	u.mockSQL.ExpectQuery("SELECT id, username, email, password, role FROM users WHERE username = $1").
		WithArgs("username").
		WillReturnError(errors.New("error"))

	_, err := u.u.GetUserByUsername("username")
	u.Error(err)
}
