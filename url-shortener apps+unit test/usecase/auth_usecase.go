package usecase

import (
	"enigmacamp.com/url-shortener/utils/service"
)

type AuthenticateUsecase interface {
	Login(username string, password string) (string, error)
}

type authenticateUsecase struct {
	userUseCase UserUseCase
	jwtService  service.JwtService
}

func (a *authenticateUsecase) Login(username string, password string) (string, error) {
	user, err := a.userUseCase.FindUserByUsernamePassword(username, password)
	if err != nil {
		return "", err
	}
	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func NewAuthenticateUsecase(userUseCase UserUseCase, jwtService service.JwtService) AuthenticateUsecase {
	return &authenticateUsecase{
		userUseCase: userUseCase,
		jwtService:  jwtService,
	}
}
