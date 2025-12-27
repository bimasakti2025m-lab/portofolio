// Mendeklarasikan unit test untuk user_usecase
package usecase

import (
	"errors"
	"testing"
	"enigmacamp.com/mini-banking/model"
	"enigmacamp.com/mini-banking/repository"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	repoMock *repository.UserRepositoryMock
	usecase  UserUseCase
}

func (s *UserUsecaseTestSuite) SetupTest() {
	s.repoMock = new(repository.UserRepositoryMock)
	s.usecase = NewUserUseCase(s.repoMock)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

func (s *UserUsecaseTestSuite) TestRegisterNewUser_Success() {
	payload := model.UserCredential{Username: "test", Password: "password", Role: "user"}
	expectedUser := model.UserCredential{Username: "test", Password: "password", Role: "user"}

	s.repoMock.On("Create", payload).Return(expectedUser, nil)

	actualUser, err := s.usecase.RegisterNewUser(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, actualUser)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestRegisterNewUser_Fail() {
	payload := model.UserCredential{Username: "test", Password: "password", Role: "user"}
	expectedError := errors.New("database error")

	s.repoMock.On("Create", payload).Return(model.UserCredential{}, expectedError)

	actualUser, err := s.usecase.RegisterNewUser(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.UserCredential{}, actualUser)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestFindUserByUsername() {
	username := "testuser"
	password := "testpass"
	expectedUser := model.UserCredential{Id: 1, Username: username, Role: "user"}

	s.repoMock.On("GetByUsernamePassword", username, password).Return(expectedUser, nil)

	actualUser, err := s.usecase.FindUserByUsernamePassword(username, password)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, actualUser)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestFindUserByUsernamePassword_NotFound() {
	username := "notfound"
	password := "wrong"
	expectedError := errors.New("user not found")

	s.repoMock.On("GetByUsernamePassword", username, password).Return(model.UserCredential{}, expectedError)

	actualUser, err := s.usecase.FindUserByUsernamePassword(username, password)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.UserCredential{}, actualUser)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestFindAllUser_Success() {
	expectedUsers := []model.UserCredential{
		{Id: 1, Username: "admin", Role: "admin"},
		{Id: 2, Username: "user1", Role: "user"},
	}

	s.repoMock.On("List").Return(expectedUsers, nil)

	actualUsers, err := s.usecase.FindAllUser()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUsers, actualUsers)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestFindAllUser_Fail() {
	expectedError := errors.New("database error")

	s.repoMock.On("List").Return(nil, expectedError)

	actualUsers, err := s.usecase.FindAllUser()

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), actualUsers)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestFindUserById_Success() {
	expectedUser := model.UserCredential{Id: 1, Username: "admin", Role: "admin"}

	s.repoMock.On("Get", uint32(1)).Return(expectedUser, nil)

	actualUser, err := s.usecase.FindUserById(1)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, actualUser)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestFindUserById_Fail() {
	expectedError := errors.New("database error")

	s.repoMock.On("Get", uint32(99)).Return(model.UserCredential{}, expectedError)

	actualUser, err := s.usecase.FindUserById(99)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.UserCredential{}, actualUser)
	s.repoMock.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestFindUserById_NotFound() {
	expectedError := errors.New("user not found")

	s.repoMock.On("Get", uint32(99)).Return(model.UserCredential{}, expectedError)

	actualUser, err := s.usecase.FindUserById(99)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.UserCredential{}, actualUser)
	s.repoMock.AssertExpectations(s.T())
}



