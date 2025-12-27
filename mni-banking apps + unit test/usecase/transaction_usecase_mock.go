package usecase

import (
	"enigmacamp.com/mini-banking/model"
	"github.com/stretchr/testify/mock"
)

type TransactionUseCaseMock struct {
	mock.Mock
}

func (m *TransactionUseCaseMock) CreateTransaction(payload model.Transaction) (model.Transaction, error) {
	args := m.Called(payload)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) ListTransaction() ([]model.Transaction, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) GetTransactionById(id uint32) (model.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) GetTransactionByUserId(userId uint32) ([]model.Transaction, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) UpdateTransaction(transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) DeleteTransaction(id uint32) error {
	args := m.Called(id)
	return args.Error(0)
}
