package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userUsecase    UserUsecase
	userRepository *repository.UserRepositoryMock
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.userRepository = new(repository.UserRepositoryMock)
	suite.userUsecase = NewUserUsecase(suite.userRepository)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

func (suite *UserUsecaseTestSuite) TestCreate_Success() {
	user := &model.User{
		Username: "testuser",
		Password: "password",
		Role:     "user",
	}

	suite.userRepository.On("Create", user).Return(user, nil)

	result, err := suite.userUsecase.Create(user)
	suite.NoError(err)
	suite.Equal(user, result)
}

func (suite *UserUsecaseTestSuite) TestCreate_Failed() {
	user := &model.User{
		Username: "testuser",
		Password: "password",
		Role:     "user",
	}

	suite.userRepository.On("Create", user).Return((*model.User)(nil), errors.New("failed to create user"))

	result, err := suite.userUsecase.Create(user)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *UserUsecaseTestSuite) TestGetAllUsers_Success() {
	users := []model.User{
		{Username: "user1", Password: "password1", Role: "user"},
		{Username: "user2", Password: "password2", Role: "admin"},
	}

	suite.userRepository.On("GetAllUsers").Return(users, nil)

	result, err := suite.userUsecase.GetAllUsers()
	suite.NoError(err)
	suite.Equal(users, result)
}

func (suite *UserUsecaseTestSuite) TestGetAllUsers_Failed() {
	suite.userRepository.On("GetAllUsers").Return([]model.User(nil), errors.New("failed to get users"))

	result, err := suite.userUsecase.GetAllUsers()
	suite.Error(err)
	suite.Nil(result)
}

func (suite *UserUsecaseTestSuite) TestGetUserByUsername_Success() {
	username := "testuser"
	user := model.User{Username: username, Password: "password", Role: "user"}

	suite.userRepository.On("GetUserByUsername", username).Return(user, nil)

	result, err := suite.userUsecase.GetUserByUsername(username)
	suite.NoError(err)
	suite.Equal(user, result)
}

func (suite *UserUsecaseTestSuite) TestGetUserByUsername_Failed() {
	username := "nonexistentuser"

	suite.userRepository.On("GetUserByUsername", username).Return(model.User{}, errors.New("user not found"))

	result, err := suite.userUsecase.GetUserByUsername(username)
	suite.Error(err)
	suite.Equal(model.User{}, result)
}




