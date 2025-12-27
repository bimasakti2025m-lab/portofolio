package usecase

import (
	"enigmacamp.com/mini-banking/model"
	"github.com/stretchr/testify/mock"
)

type AuthenticateUsecaseMock struct {
	mock.Mock
}

func (m *AuthenticateUsecaseMock) Login(username string, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *AuthenticateUsecaseMock) Register(payload model.UserCredential) (model.UserCredential, error) {
	args := m.Called(payload)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

type UserUseCaseMock struct {
	mock.Mock
}

func (m *UserUseCaseMock) RegisterNewUser(payload model.UserCredential) (model.UserCredential, error) {
	args := m.Called(payload)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (m *UserUseCaseMock) FindAllUser() ([]model.UserCredential, error) {
	args := m.Called()
	return args.Get(0).([]model.UserCredential), args.Error(1)
}

func (m *UserUseCaseMock) FindUserById(id uint32) (model.UserCredential, error) {
	args := m.Called(id)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (m *UserUseCaseMock) FindUserByUsernamePassword(username string, password string) (model.UserCredential, error) {
	args := m.Called(username, password)
	var user model.UserCredential
	// Handle nil case for the first argument if the mock returns it
	if args.Get(0) != nil {
		user = args.Get(0).(model.UserCredential)
	}
	return user, args.Error(1)
}
