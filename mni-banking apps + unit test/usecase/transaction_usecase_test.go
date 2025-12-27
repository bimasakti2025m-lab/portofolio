package usecase

import (
	"errors"
	"testing"

	"enigmacamp.com/mini-banking/model"
	"enigmacamp.com/mini-banking/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionUsecaseTestSuite struct {
	suite.Suite
	repoMock *repository.TransactionRepositoryMock
	usecase  TransactionUseCase
}

func (s *TransactionUsecaseTestSuite) SetupTest() {
	s.repoMock = new(repository.TransactionRepositoryMock)
	s.usecase = NewTransactionUseCase(s.repoMock)
}

func TestTransactionUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionUsecaseTestSuite))
}

func (s *TransactionUsecaseTestSuite) TestCreateTransaction_Success() {
	payload := model.Transaction{FromUserID: 1, ToUserID: 2, Amount: 100}
	expected := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100}
	s.repoMock.On("Create", payload).Return(expected, nil)

	actual, err := s.usecase.CreateTransaction(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestCreateTransaction_Fail() {
	payload := model.Transaction{FromUserID: 1, ToUserID: 2, Amount: 100}
	expectedError := errors.New("failed to create")
	s.repoMock.On("Create", payload).Return(model.Transaction{}, expectedError)

	actual, err := s.usecase.CreateTransaction(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Transaction{}, actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestListTransaction_Success() {
	expected := []model.Transaction{{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100}}
	s.repoMock.On("List").Return(expected, nil)

	actual, err := s.usecase.ListTransaction()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestListTransaction_Fail() {
	expectedError := errors.New("failed to list")
	s.repoMock.On("List").Return(nil, expectedError)

	actual, err := s.usecase.ListTransaction()

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestGetTransactionById_Success() {
	id := uint32(1)
	expected := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100}
	s.repoMock.On("Get", id).Return(expected, nil)

	actual, err := s.usecase.GetTransactionById(id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestGetTransactionById_Fail() {
	id := uint32(1)
	expectedError := errors.New("not found")
	s.repoMock.On("Get", id).Return(model.Transaction{}, expectedError)

	actual, err := s.usecase.GetTransactionById(id)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Transaction{}, actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestGetTransactionByUserId_Success() {
	userId := uint32(1)
	expected := []model.Transaction{{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100}}
	s.repoMock.On("GetByUserId", userId).Return(expected, nil)

	actual, err := s.usecase.GetTransactionByUserId(userId)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestGetTransactionByUserId_Fail() {
	userId := uint32(1)
	expectedError := errors.New("failed to get by user id")
	s.repoMock.On("GetByUserId", userId).Return(nil, expectedError)

	actual, err := s.usecase.GetTransactionByUserId(userId)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Nil(s.T(), actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestUpdateTransaction_Success() {
	payload := model.Transaction{ID: 1, Amount: 200}
	s.repoMock.On("Update", payload).Return(payload, nil)

	actual, err := s.usecase.UpdateTransaction(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), payload, actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestUpdateTransaction_Fail() {
	payload := model.Transaction{ID: 1, Amount: 200}
	expectedError := errors.New("failed to update")
	s.repoMock.On("Update", payload).Return(model.Transaction{}, expectedError)

	actual, err := s.usecase.UpdateTransaction(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Transaction{}, actual)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestDeleteTransaction_Success() {
	id := uint32(1)
	s.repoMock.On("Delete", id).Return(nil)
	err := s.usecase.DeleteTransaction(id)
	assert.NoError(s.T(), err)
	s.repoMock.AssertExpectations(s.T())
}

func (s *TransactionUsecaseTestSuite) TestDeleteTransaction_Fail() {
	id := uint32(1)
	expectedError := errors.New("failed to delete")
	s.repoMock.On("Delete", id).Return(expectedError)
	err := s.usecase.DeleteTransaction(id)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	s.repoMock.AssertExpectations(s.T())
}


