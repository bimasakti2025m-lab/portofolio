package controller_mock

import (
	"basic-JWT/model"

	"github.com/stretchr/testify/mock"
)

type AuthenticationUsecaseMock struct {
	mock.Mock
}

func (a *AuthenticationUsecaseMock) Login(username string, password string) (string, error) {
	args := a.Called(username, password)
	return args.String(0), args.Error(1)
}

func (a *AuthenticationUsecaseMock) Register(username string, password string) (model.User, error) {
	args := a.Called(username, password)
	return args.Get(0).(model.User), args.Error(1)
}
