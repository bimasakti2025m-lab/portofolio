package repository

import (
	"enigmacamp.com/url-shortener/model"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user model.UserCredential) (model.UserCredential, error) {
	args := m.Called(user)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (m *UserRepositoryMock) List() ([]model.UserCredential, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.UserCredential), args.Error(1)
}

func (m *UserRepositoryMock) Get(id uint32) (model.UserCredential, error) {
	args := m.Called(id)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (m *UserRepositoryMock) GetByUsernamePassword(username string, password string) (model.UserCredential, error) {
	args := m.Called(username, password)
	return args.Get(0).(model.UserCredential), args.Error(1)
}
