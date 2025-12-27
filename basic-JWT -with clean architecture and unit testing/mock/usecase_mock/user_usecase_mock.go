package usecase_mock

import (
	"basic-JWT/model"

	"github.com/stretchr/testify/mock"
)

type UserUseCaseMock struct {
	mock.Mock
}

func (u *UserUseCaseMock) Create(user *model.User) (*model.User, error) {
	args := u.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)

}

func (u *UserUseCaseMock) GetAllUsers() ([]model.User, error) {
	args := u.Called()
	return args.Get(0).([]model.User), args.Error(1)

}

func (u *UserUseCaseMock) GetUserByUsername(username string) (model.User, error) {
	args := u.Called(username)
	return args.Get(0).(model.User), args.Error(1)

}
