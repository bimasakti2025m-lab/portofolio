package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/utils/service"

	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationUsecase interface {
	Register(username string, email, password, role string) (model.User, error)
	Login(username string, password string) (string, error)
}

type authenticationUsecase struct {
	userUsecase UserUsecase
	jwtService  service.JWTservice
}

func (au *authenticationUsecase) Login(username string, password string) (string, error) {
	user, err := au.userUsecase.GetUserByUsername(username)
	if err != nil {
		return "Failed to get user by username", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "Failed to compare password", err
	}

	token := au.jwtService.CreateToken(user)

	return token, nil

}

func (au *authenticationUsecase) Register(username, email, password, role string) (model.User, error) {
	// Check if username is taken
	user, err := au.userUsecase.GetUserByUsername(username)
	if err == nil {
		return model.User{}, fmt.Errorf("username '%s' is already taken", username)
	}

	// Check if username is empty
	if username == "" {
		return model.User{}, fmt.Errorf("username cannot be empty")
	}

	// Check if email is empty
	if email == "" {
		return model.User{}, fmt.Errorf("email cannot be empty")
	}

	// Check if password is empty
	if password == "" {
		return model.User{}, fmt.Errorf("password cannot be empty")
	}

	// Check if role is empty
	if role == "" {
		return model.User{}, fmt.Errorf("role cannot be empty")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("Failed to hashing password")
	}

	password = string(hashedPassword)

	// buat role bisa di isi admin atau user
	user = model.User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     role,
	}

	// Create the user in the database
	_, err = au.userUsecase.Create(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func NewAuthenticationUsecase(userUsecase UserUsecase, jwtService service.JWTservice) AuthenticationUsecase {
	return &authenticationUsecase{
		userUsecase: userUsecase,
		jwtService:  jwtService,
	}
}
