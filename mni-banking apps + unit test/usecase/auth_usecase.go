package usecase

import (
	"fmt"

	"enigmacamp.com/mini-banking/model"
	"enigmacamp.com/mini-banking/utils/service"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateUsecase interface {
	Login(username string, password string) (string, error)
	Register(payload model.UserCredential) (model.UserCredential, error)
}

type authenticateUsecase struct {
	userUseCase UserUseCase
	jwtService  service.JwtService
}

func (a *authenticateUsecase) Login(username string, password string) (string, error) {
	// TODO : GET USER BY USERNAME AND PASSWORD
	user, err := a.userUseCase.FindUserByUsernamePassword(username, password)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	token, err := a.jwtService.CreateToken(user)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return token, nil
}

func (a *authenticateUsecase) Register(payload model.UserCredential) (model.UserCredential, error) {
	// 1. Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserCredential{}, fmt.Errorf("failed to hash password: %w", err)
	}
	payload.Password = string(hashedPassword)

	// 3. Simpan pengguna ke database
	return a.userUseCase.RegisterNewUser(payload)
}

func NewAuthenticateUsecase(userUseCase UserUseCase, jwtService service.JwtService) AuthenticateUsecase {
	return &authenticateUsecase{
		userUseCase: userUseCase,
		jwtService:  jwtService,
	}
}
