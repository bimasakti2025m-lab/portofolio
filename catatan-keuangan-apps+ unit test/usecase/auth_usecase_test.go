package usecase

import (
	"fmt"
	"testing"

	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/entity/dto"
	"enigmacamp.com/livecode-catatan-keuangan/mock/service_mock"
	"enigmacamp.com/livecode-catatan-keuangan/mock/usecase_mock"
	"github.com/stretchr/testify/suite"
)

type authUCSuite struct {
	suite.Suite
	userUseCase *usecase_mock.UserUsecaseMock
	jwtService  *service_mock.JwtServiceMock
	authUC      AuthUseCase
}

func TestAuthUCSuite(t *testing.T) {
	suite.Run(t, new(authUCSuite))
}

func (a *authUCSuite) SetupTest() {
	a.userUseCase = new(usecase_mock.UserUsecaseMock)
	a.jwtService = new(service_mock.JwtServiceMock)
	a.authUC = NewAuthUseCase(a.userUseCase, a.jwtService)
}

func (a *authUCSuite) TestRegister_failed() {
	newUser := dto.AuthRequestDto{
		Username: "failed",
		Password: "password failed",
	}

	a.userUseCase.On("RegisterNewUser", entity.User{Username: newUser.Username, Password: newUser.Password}).Return(entity.User{}, fmt.Errorf("failed")).Once()

	user, err := a.authUC.Register(newUser)
	a.Equal(user, entity.User{})
	a.NotNil(err)
}

func (a *authUCSuite) TestRegister_success() {
	newUser := dto.AuthRequestDto{
		Username: "success",
		Password: "password success",
	}

	a.userUseCase.On("RegisterNewUser", entity.User{Username: newUser.Username, Password: newUser.Password}).Return(entity.User{
		Username: newUser.Username,
		Password: newUser.Password,
	}, nil).Once()

	user, err := a.authUC.Register(newUser)
	a.Equal(user, entity.User{
		Username: newUser.Username,
		Password: newUser.Password,
	})
	a.Nil(err)
}
