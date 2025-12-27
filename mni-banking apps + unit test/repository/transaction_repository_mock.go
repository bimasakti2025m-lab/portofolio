package repository

import (
	"enigmacamp.com/mini-banking/model"
	"github.com/stretchr/testify/mock"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func (m *TransactionRepositoryMock) Create(transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionRepositoryMock) List() ([]model.Transaction, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *TransactionRepositoryMock) Get(id uint32) (model.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionRepositoryMock) GetByUserId(userId uint32) ([]model.Transaction, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *TransactionRepositoryMock) Update(transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionRepositoryMock) Delete(id uint32) error {
	args := m.Called(id)
	return args.Error(0)
}
