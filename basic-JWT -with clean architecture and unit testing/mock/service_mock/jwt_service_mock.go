package service_mock

import (
	"basic-JWT/model"
	modelutils "basic-JWT/utils/model_utils"

	"github.com/stretchr/testify/mock"
)

type JWTServiceMock struct {
	mock.Mock
}

func (j *JWTServiceMock) CreateToken(user model.User) string{
	args := j.Called(user)
	return args.String(0)
}
func  (j *JWTServiceMock) VerifyToken(tokenString string) (*modelutils.JwtPayloadClaims, error){
	args := j.Called(tokenString)
	return args.Get(0).(*modelutils.JwtPayloadClaims), args.Error(1)
}