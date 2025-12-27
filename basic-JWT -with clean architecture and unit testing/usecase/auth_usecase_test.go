package usecase_test

import (
	"basic-JWT/mock/service_mock"
	"basic-JWT/mock/usecase_mock"
	"basic-JWT/model"
	"basic-JWT/usecase"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type authUCSuite struct {
	suite.Suite
	authUC      usecase.AuthenticationUsecase
	jwtService  *service_mock.JWTServiceMock
	UserUsecase *usecase_mock.UserUseCaseMock
}

func TestAuthUcSuite(t *testing.T) {
	suite.Run(t, new(authUCSuite))
}

func (a *authUCSuite) SetupTest() {
	a.UserUsecase = new(usecase_mock.UserUseCaseMock)
	a.jwtService = new(service_mock.JWTServiceMock)
	a.authUC = usecase.NewAuthenticationUsecase(a.UserUsecase, a.jwtService)
}

func (a *authUCSuite) TestLogin() {
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	a.UserUsecase.On("GetUserByUsername", username).Return(user, nil)
	a.jwtService.On("CreateToken", user).Return("token")

	token, err := a.authUC.Login(username, password)
	a.NoError(err)
	a.Equal("token", token)
}

func (a *authUCSuite) TestLogin_UserNotFound() {
	username := "username"
	password := "password"

	a.UserUsecase.On("GetUserByUsername", username).Return(model.User{}, nil)

	_, err := a.authUC.Login(username, password)
	a.Error(err)
}

func (a *authUCSuite) TestLogin_PasswordMismatch() {
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	a.UserUsecase.On("GetUserByUsername", username).Return(user, nil)

	_, err := a.authUC.Login(username, "wrongpassword")
	a.Error(err)
}

func (a *authUCSuite) TestRegister_UsernameTaken() {
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	a.UserUsecase.On("GetUserByUsername", username).Return(user, nil)

	_, err := a.authUC.Register(username, password)
	a.Error(err)
}

func (a *authUCSuite) TestRegister_Failed() {
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	a.UserUsecase.On("GetUserByUsername", username).Return(user, nil)
	a.UserUsecase.On("Create", &model.User{Username: username, Password: string(hashedPassword), Role: "user"}).Return(nil, errors.New("error creating user"))

	_, err := a.authUC.Register(username, password)
	a.Error(err)
}

func (a *authUCSuite) TestRegister_EmptyUsername() {
	username := ""
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	a.UserUsecase.On("GetUserByUsername", username).Return(user, errors.New("username cannot be empty"))
	a.UserUsecase.On("Create", &model.User{Username: username, Password: string(hashedPassword), Role: "user"}).Return(nil, errors.New("username cannot be empty"))

	_, err := a.authUC.Register(username, password)
	a.Error(err)
}

func (a *authUCSuite) TestRegister_EmptyPassword() {
	username := "username"
	password := ""
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	a.UserUsecase.On("GetUserByUsername", username).Return(user, errors.New("password cannot be empty"))
	a.UserUsecase.On("Create", &model.User{Username: username, Password: string(hashedPassword), Role: "user"}).Return(nil, errors.New("password cannot be empty"))

	_, err := a.authUC.Register(username, password)
	a.Error(err)
}
