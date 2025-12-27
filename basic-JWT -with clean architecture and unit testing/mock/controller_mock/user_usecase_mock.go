package controller_mock

import (
	"basic-JWT/model"

	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (u *UserUsecaseMock) Create(user *model.User) (*model.User, error){
	args := u.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *UserUsecaseMock) GetAllUsers() ([]model.User, error) {
	args := u.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (u *UserUsecaseMock) GetUserByUsername(username string) (model.User, error) {
	args := u.Called(username)
	return args.Get(0).(model.User), args.Error(1)
}