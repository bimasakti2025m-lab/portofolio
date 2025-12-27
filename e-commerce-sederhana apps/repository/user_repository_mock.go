package repository

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

// buat struct mock repository
type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepositoryMock) GetAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByUsername(username string) (model.User, error) {
	args := m.Called(username)
	return args.Get(0).(model.User), args.Error(1)
}
