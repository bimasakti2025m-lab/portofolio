package usecase

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type UserUseCaseMock struct {
	mock.Mock
}

func (m *UserUseCaseMock) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserUseCaseMock) GetAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *UserUseCaseMock) GetUserByUsername(username string) (model.User, error) {
	args := m.Called(username)
	return args.Get(0).(model.User), args.Error(1)
}