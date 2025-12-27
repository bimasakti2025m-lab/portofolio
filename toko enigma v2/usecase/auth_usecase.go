package usecase

import (
	"fmt"

	"enigmacamp.com/toko-enigma/model"
	"enigmacamp.com/toko-enigma/utils/service"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateUsecase interface {
	Login(username string, password string) (string, error)
	// add fitur register
	Register(payload model.UserCredential) (model.UserCredential, error)
}

type authenticateUsecase struct {
	userUseCase UserUseCase
	jwtService  service.JwtService
}

func (a *authenticateUsecase) Login(username string, password string) (string, error) {
	// 1. Ambil user berdasarkan username saja
	user, err := a.userUseCase.FindUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// 2. Bandingkan password yang diinput dengan hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}
	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *authenticateUsecase) Register(payload model.UserCredential) (model.UserCredential, error) {
	// Tambahan: Periksa apakah username sudah ada
	_, err := a.userUseCase.FindUserByUsername(payload.Username)
	if err == nil { // Jika tidak ada error, berarti user ditemukan (sudah ada)
		return model.UserCredential{}, fmt.Errorf("username '%s' already exists", payload.Username)
	}

	hasshedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserCredential{}, err
	}
	payload.Password = string(hasshedPassword)

	user, err := a.userUseCase.RegisterNewUser(payload) // Pastikan ini memanggil repo.Create
	if err != nil {
		return model.UserCredential{}, err
	}

	return user, nil
}

func NewAuthenticateUsecase(userUseCase UserUseCase, jwtService service.JwtService) AuthenticateUsecase {
	return &authenticateUsecase{
		userUseCase: userUseCase,
		jwtService:  jwtService,
	}
}
