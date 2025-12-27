package usecase

import (
	"errors"
	"testing"

	"enigmacamp.com/toko-enigma/model"
	"enigmacamp.com/toko-enigma/utils/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateUsecaseTestSuite struct {
	suite.Suite
	userUseCaseMock *UserUseCaseMock
	jwtServiceMock  *service.JwtServiceMock
	usecase         AuthenticateUsecase
}

func (s *AuthenticateUsecaseTestSuite) SetupTest() {
	s.userUseCaseMock = new(UserUseCaseMock)
	s.jwtServiceMock = new(service.JwtServiceMock)
	s.usecase = NewAuthenticateUsecase(s.userUseCaseMock, s.jwtServiceMock)
}

func TestAuthenticateUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticateUsecaseTestSuite))
}

func (s *AuthenticateUsecaseTestSuite) TestLogin_Success() {
	username := "testuser"
	password := "testpass"
	// Hash the password for the expected user to simulate the stored password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	expectedUser := model.UserCredential{Id: 1, Username: username, Password: string(hashedPassword), Role: "user"}
	expectedToken := "mocked_jwt_token"

	s.userUseCaseMock.On("FindUserByUsername", username).Return(expectedUser, nil)

	s.userUseCaseMock.On("FindUserByUsernamePassword", username, password).Return(expectedUser, nil)
	s.jwtServiceMock.On("CreateToken", expectedUser).Return(expectedToken, nil)

	token, err := s.usecase.Login(username, password)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedToken, token)
	
	s.jwtServiceMock.AssertExpectations(s.T())
}

func (s *AuthenticateUsecaseTestSuite) TestLogin_InvalidCredentials() {
	username := "testuser"
	password := "wrongpass"
	// Simulate a user found by username, but with a different password
	expectedUser := model.UserCredential{Id: 1, Username: username, Password: "hashed_correct_password", Role: "user"}

	s.userUseCaseMock.On("FindUserByUsername", username).Return(expectedUser, nil)

	token, err := s.usecase.Login(username, password)

	assert.Error(s.T(), err)
	assert.Empty(s.T(), token)
	assert.Equal(s.T(), "invalid username or password", err.Error())
	s.userUseCaseMock.AssertExpectations(s.T())
}

func (s *AuthenticateUsecaseTestSuite) TestLogin_TokenCreationFail() {
	username := "testuser"
	password := "testpass"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	expectedUser := model.UserCredential{Id: 1, Username: username, Password: string(hashedPassword), Role: "user"}
	expectedError := errors.New("failed to create token")
	s.userUseCaseMock.On("FindUserByUsername", username).Return(expectedUser, nil)

	s.userUseCaseMock.On("FindUserByUsernamePassword", username, password).Return(expectedUser, nil)
	s.jwtServiceMock.On("CreateToken", expectedUser).Return("", expectedError)

	token, err := s.usecase.Login(username, password)

	assert.Error(s.T(), err)
	assert.Empty(s.T(), token)
	assert.Equal(s.T(), expectedError, err)
	
	s.jwtServiceMock.AssertExpectations(s.T())
}
func (s *AuthenticateUsecaseTestSuite) TestLogin_UserNotFound() {
	username := "nonexistent"
	password := "anypass"
	expectedError := errors.New("failed to get user by username") // Error from repository

	s.userUseCaseMock.On("FindUserByUsername", username).Return(model.UserCredential{}, expectedError)

	token, err := s.usecase.Login(username, password)

	assert.Error(s.T(), err)
	assert.Empty(s.T(), token)
	assert.Equal(s.T(), "invalid username or password", err.Error()) // Usecase should return generic error
	s.userUseCaseMock.AssertExpectations(s.T())
	s.jwtServiceMock.AssertNotCalled(s.T(), "CreateToken", mock.Anything)
}

func (s *AuthenticateUsecaseTestSuite) TestRegister_Success() {
	payload := model.UserCredential{Username: "newuser", Password: "newpass", Role: "user"}
	expectedUser := model.UserCredential{Id: 2, Username: "newuser", Role: "user"}

	// Mock FindUserByUsername to return an error (user not found)
	s.userUseCaseMock.On("FindUserByUsername", payload.Username).Return(model.UserCredential{}, errors.New("user not found"))
	// Mock RegisterNewUser to return the created user
	s.userUseCaseMock.On("RegisterNewUser", mock.AnythingOfType("model.UserCredential")).Return(expectedUser, nil)

	actualUser, err := s.usecase.Register(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser.Id, actualUser.Id)
	assert.Equal(s.T(), expectedUser.Username, actualUser.Username)
	assert.Equal(s.T(), expectedUser.Role, actualUser.Role)
	s.userUseCaseMock.AssertExpectations(s.T())
}

func (s *AuthenticateUsecaseTestSuite) TestRegister_UserAlreadyExists() {
	payload := model.UserCredential{Username: "existinguser", Password: "newpass", Role: "user"}
	existingUser := model.UserCredential{Id: 1, Username: "existinguser", Role: "user"}

	// Mock FindUserByUsername to return an existing user (user found)
	s.userUseCaseMock.On("FindUserByUsername", payload.Username).Return(existingUser, nil)

	actualUser, err := s.usecase.Register(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "username 'existinguser' already exists", err.Error())
	assert.Equal(s.T(), model.UserCredential{}, actualUser)
	s.userUseCaseMock.AssertExpectations(s.T())
	// Ensure RegisterNewUser was not called
	s.userUseCaseMock.AssertNotCalled(s.T(), "RegisterNewUser", mock.Anything)
}
