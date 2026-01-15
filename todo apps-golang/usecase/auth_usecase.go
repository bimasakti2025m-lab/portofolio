package usecase

import (
	"fmt"

	"enigmacamp.com/golang-jwt/model"
	"enigmacamp.com/golang-jwt/utils/service"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateUsecase interface {
	Login(payload model.UserCredential) (string, error)
	Register(payload model.UserCredential) (model.UserCredential, error)
}

type authenticateUsecase struct {
	userUseCase UserUseCase
	jwtService  service.JwtService
}

func (a *authenticateUsecase) Login(payload model.UserCredential) (string, error) {
	user, err := a.userUseCase.FindUserByUsername(payload.Username)
	if err != nil {
		return "", fmt.Errorf("invalid username")
	}

	// Verifikasi password yang di-hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Mendeklarasikan fitur register
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
