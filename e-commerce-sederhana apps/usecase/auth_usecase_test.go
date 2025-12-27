package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/utils/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseTestSuite struct {
	suite.Suite
	userUseCase *UserUseCaseMock
	jwtService  *service.JWTServiceMock
	usecase     AuthenticationUsecase
}

func (suite *AuthUsecaseTestSuite) SetupTest() {
	suite.userUseCase = new(UserUseCaseMock)
	suite.jwtService = new(service.JWTServiceMock)
	suite.usecase = NewAuthenticationUsecase(suite.userUseCase, suite.jwtService)
}

func TestAuthUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUsecaseTestSuite))
}

func (suite *AuthUsecaseTestSuite) TestRegister_Success() {
	username := "newuser"
	email := "new@mail.com"
	password := "password"
	role := "user"

	// Expect GetUserByUsername to return error (user not found)
	suite.userUseCase.On("GetUserByUsername", username).Return(model.User{}, errors.New("user not found"))

	// Expect Create to be called with hashed password
	suite.userUseCase.On("Create", mock.MatchedBy(func(u *model.User) bool {
		return u.Username == username && u.Email == email && u.Role == role && u.Password != password
	})).Return(&model.User{Username: username}, nil)

	user, err := suite.usecase.Register(username, email, password, role)
	suite.NoError(err)
	suite.Equal(username, user.Username)
}

func (suite *AuthUsecaseTestSuite) TestRegister_UserExists() {
	username := "existinguser"
	suite.userUseCase.On("GetUserByUsername", username).Return(model.User{Username: username}, nil)

	_, err := suite.usecase.Register(username, "mail", "pass", "role")
	suite.Error(err)
	suite.Contains(err.Error(), "already taken")
}

func (suite *AuthUsecaseTestSuite) TestLogin_Success() {
	username := "user"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	suite.userUseCase.On("GetUserByUsername", username).Return(user, nil)
	suite.jwtService.On("CreateToken", user).Return("mocked_token")

	token, err := suite.usecase.Login(username, password)
	suite.NoError(err)
	suite.Equal("mocked_token", token)
}

func (suite *AuthUsecaseTestSuite) TestLogin_WrongPassword() {
	username := "user"
	password := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	user := model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	suite.userUseCase.On("GetUserByUsername", username).Return(user, nil)

	token, err := suite.usecase.Login(username, password)

	suite.Error(err)
	suite.Contains(token, "Failed to compare password")

}

func (suite *AuthUsecaseTestSuite) TestLogin_UserNotFound() {
	suite.userUseCase.On("GetUserByUsername", "unknown").Return(model.User{}, errors.New("not found"))
	_, err := suite.usecase.Login("unknown", "pass")
	suite.Error(err)
}
