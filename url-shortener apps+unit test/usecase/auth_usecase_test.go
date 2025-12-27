package usecase

import (
	"errors"
	"testing"

	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/utils/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	expectedUser := model.UserCredential{Id: 1, Username: username, Role: "user"}
	expectedToken := "mocked_jwt_token"

	s.userUseCaseMock.On("FindUserByUsernamePassword", username, password).Return(expectedUser, nil)
	s.jwtServiceMock.On("CreateToken", expectedUser).Return(expectedToken, nil)

	token, err := s.usecase.Login(username, password)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedToken, token)
	s.userUseCaseMock.AssertExpectations(s.T())
	s.jwtServiceMock.AssertExpectations(s.T())
}

func (s *AuthenticateUsecaseTestSuite) TestLogin_InvalidCredentials() {
	username := "testuser"
	password := "wrongpass"
	expectedError := errors.New("invalid username or password")

	s.userUseCaseMock.On("FindUserByUsernamePassword", username, password).Return(model.UserCredential{}, expectedError)

	token, err := s.usecase.Login(username, password)

	assert.Error(s.T(), err)
	assert.Empty(s.T(), token)
	assert.Equal(s.T(), expectedError, err)
	s.userUseCaseMock.AssertExpectations(s.T())
}

func (s *AuthenticateUsecaseTestSuite) TestLogin_TokenCreationFail() {
	username := "testuser"
	password := "testpass"
	expectedUser := model.UserCredential{Id: 1, Username: username, Role: "user"}
	expectedError := errors.New("failed to create token")

	s.userUseCaseMock.On("FindUserByUsernamePassword", username, password).Return(expectedUser, nil)
	s.jwtServiceMock.On("CreateToken", expectedUser).Return("", expectedError)

	token, err := s.usecase.Login(username, password)

	assert.Error(s.T(), err)
	assert.Empty(s.T(), token)
	assert.Equal(s.T(), expectedError, err)
	s.userUseCaseMock.AssertExpectations(s.T())
	s.jwtServiceMock.AssertExpectations(s.T())
}
