package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/mini-banking/middleware"
	"enigmacamp.com/mini-banking/model"
	"enigmacamp.com/mini-banking/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	usecaseMock *usecase.TransactionUseCaseMock
	authMid     middleware.AuthMiddleware
}

func (s *TransactionControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.usecaseMock = new(usecase.TransactionUseCaseMock)
	s.authMid = &MockAuthMiddleware{} // Gunakan mock middleware
}

func TestTransactionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionControllerTestSuite))
}

func (s *TransactionControllerTestSuite) TestCreateTransaction_Success() {
	payload := model.Transaction{FromUserID: 1, ToUserID: 2, Amount: 100}
	expected := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100}
	s.usecaseMock.On("CreateTransaction", payload).Return(expected, nil)

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actual model.Transaction
	json.Unmarshal(w.Body.Bytes(), &actual)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	assert.Equal(s.T(), expected, actual)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestCreateTransaction_BindingError() {
	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *TransactionControllerTestSuite) TestCreateTransaction_UsecaseError() {
	payload := model.Transaction{FromUserID: 1, ToUserID: 2, Amount: 100}
	s.usecaseMock.On("CreateTransaction", payload).Return(model.Transaction{}, errors.New("some error"))

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestListTransaction_Success() {
	expected := []model.Transaction{{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100}}
	s.usecaseMock.On("ListTransaction").Return(expected, nil)

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actual []model.Transaction
	json.Unmarshal(w.Body.Bytes(), &actual)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expected, actual)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestListTransaction_Empty() {
	s.usecaseMock.On("ListTransaction").Return([]model.Transaction{}, nil)

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Contains(s.T(), w.Body.String(), "List transaction empty")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestListTransaction_UsecaseError() {
	s.usecaseMock.On("ListTransaction").Return(nil, errors.New("some error"))

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestGetTransactionById_Success() {
	expected := model.Transaction{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100}
	s.usecaseMock.On("GetTransactionById", uint32(1)).Return(expected, nil)

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actual model.Transaction
	json.Unmarshal(w.Body.Bytes(), &actual)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expected, actual)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestGetTransactionById_UsecaseError() {
	s.usecaseMock.On("GetTransactionById", uint32(1)).Return(model.Transaction{}, errors.New("not found"))

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestGetTransactionByUserId_Success() {
	expected := []model.Transaction{{ID: 1, FromUserID: 1, ToUserID: 2, Amount: 100}}
	s.usecaseMock.On("GetTransactionByUserId", uint32(1)).Return(expected, nil)

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions/user/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actual []model.Transaction
	json.Unmarshal(w.Body.Bytes(), &actual)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), expected, actual)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestUpdateTransaction_Success() {
	payload := model.Transaction{ID: 1, Amount: 200}
	s.usecaseMock.On("UpdateTransaction", payload).Return(payload, nil)

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/transactions", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var actual model.Transaction
	json.Unmarshal(w.Body.Bytes(), &actual)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), payload, actual)
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestDeleteTransaction_Success() {
	s.usecaseMock.On("DeleteTransaction", uint32(1)).Return(nil)

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/transactions/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Contains(s.T(), w.Body.String(), "Transaction deleted successfully")
	s.usecaseMock.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestDeleteTransaction_UsecaseError() {
	s.usecaseMock.On("DeleteTransaction", uint32(1)).Return(errors.New("some error"))

	NewTransactionController(s.usecaseMock, s.router.Group("/api/v1"), s.authMid).Route()

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/transactions/1", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	s.usecaseMock.AssertExpectations(s.T())
}
