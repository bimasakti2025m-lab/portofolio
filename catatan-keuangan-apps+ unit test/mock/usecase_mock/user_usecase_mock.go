package usecase_mock

import (
	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (u *UserUsecaseMock) Create(payload entity.User) (entity.User, error) {
	args := u.Called(payload)
	return args.Get(0).(entity.User), args.Error(1)
}

func (u *UserUsecaseMock) Get(id string) (entity.User, error) {
	args := u.Called(id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (u *UserUsecaseMock) GetByUsername(username string) (entity.User, error) {
	args := u.Called(username)
	return args.Get(0).(entity.User), args.Error(1)
}

func (u *UserUsecaseMock) RegisterNewUser(payload entity.User) (entity.User, error) {
	args := u.Called(payload)
	return args.Get(0).(entity.User), args.Error(1)
}

func (u *UserUsecaseMock) FindUserByID(id string) (entity.User, error) {
	args := u.Called(id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (u *UserUsecaseMock) FindUserByUsernamePassword(username, password string) (entity.User, error) {
	args := u.Called(username, password)
	return args.Get(0).(entity.User), args.Error(1)
}
