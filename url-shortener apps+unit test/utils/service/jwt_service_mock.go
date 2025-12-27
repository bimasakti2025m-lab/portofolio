package service

import (
	"enigmacamp.com/url-shortener/model"
	modelutils "enigmacamp.com/url-shortener/utils/model_utils"
	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func (m *JwtServiceMock) CreateToken(user model.UserCredential) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *JwtServiceMock) VerifyToken(token string) (*modelutils.JwtPayloadClaims, error) {
	args := m.Called(token)
	return args.Get(0).(*modelutils.JwtPayloadClaims), args.Error(1)
}


