package service

import (
	"E-commerce-Sederhana/model"
	modelutils "E-commerce-Sederhana/utils/model_utils"

	"github.com/stretchr/testify/mock"
)

type JWTServiceMock struct {
	mock.Mock
}

func (m *JWTServiceMock) CreateToken(user model.User) string {
	args := m.Called(user)
	return args.String(0)
}

func (m *JWTServiceMock) VerifyToken(token string) (*modelutils.JwtPayloadClaims, error) {
	args := m.Called(token)
	return args.Get(0).(*modelutils.JwtPayloadClaims), args.Error(1)
}