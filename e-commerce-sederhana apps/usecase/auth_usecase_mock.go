package usecase

import (
	"E-commerce-Sederhana/model"

	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) Register(username, email, password, role string) (model.User, error) {
	args := m.Called(username, email, password, role)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *AuthUseCaseMock) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}
