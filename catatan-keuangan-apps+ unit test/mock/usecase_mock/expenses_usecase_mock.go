package usecase_mock

import (
	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/shared/model"
	"github.com/stretchr/testify/mock"
)

type ExpensesUsecaseMock struct {
	mock.Mock
}

func (e *ExpensesUsecaseMock) Create(payload entity.Expense) (entity.Expense, error) {
	args := e.Called(payload)
	return args.Get(0).(entity.Expense), args.Error(1)
}

func (e *ExpensesUsecaseMock) Get(id string) (entity.Expense, error) {
	args := e.Called(id)
	return args.Get(0).(entity.Expense), args.Error(1)
}

func (e *ExpensesUsecaseMock) GetByTransaction(transactionType string, user string) ([]entity.Expense, error) {
	args := e.Called(transactionType, user)
	return args.Get(0).([]entity.Expense), args.Error(1)
}

func (e *ExpensesUsecaseMock) GetBalance(user string) (float64, error) {
	args := e.Called(user)
	return args.Get(0).(float64), args.Error(1)
}

func(e *ExpensesUsecaseMock) List(page, size int, startDate, endDate string, user string) ([]entity.Expense, model.Paging, error) {
	args := e.Called(page, size, startDate, endDate, user)
	return args.Get(0).([]entity.Expense), args.Get(1).(model.Paging), args.Error(2)
}


func (e *ExpensesUsecaseMock) RegisterNewExpense(payload entity.Expense) (entity.Expense, error) {
	args := e.Called(payload)
	return args.Get(0).(entity.Expense), args.Error(1)
}

func (e *ExpensesUsecaseMock) FindAllExpense(page, size int, startDate, endDate string, user string) ([]entity.Expense, model.Paging, error) {
	args := e.Called(page, size, startDate, endDate, user)
	return args.Get(0).([]entity.Expense), args.Get(1).(model.Paging), args.Error(2)
}

func (e *ExpensesUsecaseMock) FindExpenseByID(id string) (entity.Expense, error) {
	args := e.Called(id)
	return args.Get(0).(entity.Expense), args.Error(1)
}

func (e *ExpensesUsecaseMock) FindExpenseByTransactionType(transactionType string, user string) ([]entity.Expense, error) {
	args := e.Called(transactionType, user)
	return args.Get(0).([]entity.Expense), args.Error(1)
}

