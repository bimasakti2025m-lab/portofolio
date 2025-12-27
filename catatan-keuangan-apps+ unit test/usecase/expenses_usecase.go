package usecase

import (
	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/repository"
	"enigmacamp.com/livecode-catatan-keuangan/shared/model"
	"fmt"
	"strings"
	"time"
)

type ExpenseUseCase interface {
	RegisterNewExpense(payload entity.Expense) (entity.Expense, error)
	FindAllExpense(page, size int, startDate, endDate string, user string) ([]entity.Expense, model.Paging, error)
	FindExpenseByID(id string) (entity.Expense, error)
	FindExpenseByTransactionType(transactionType string, user string) ([]entity.Expense, error)
}

type expenseUseCase struct {
	repo repository.ExpenseRepository
}

func (e *expenseUseCase) RegisterNewExpense(payload entity.Expense) (entity.Expense, error) {
	if !payload.IsRequiredFields() {
		return entity.Expense{}, fmt.Errorf("opps, required fields")
	}

	if !payload.IsTransactionTypeValid() {
		return entity.Expense{}, fmt.Errorf("opps, transaction type must CREDIT or DEBIT")
	}

	balance, _ := e.repo.GetBalance(payload.UserId)
	if balance < 0 {
		return entity.Expense{}, fmt.Errorf("opps, error get balance")
	}

	if payload.TransactionType == "CREDIT" {
		payload.Balance = balance + payload.Amount
	} else {
		if payload.Amount > balance {
			return entity.Expense{}, fmt.Errorf("opps, balance not enough")
		}
		payload.Balance = balance - payload.Amount
	}

	payload.Date = time.Now()
	payload.UpdatedAt = time.Now()
	return e.repo.Create(payload)
}

func (e *expenseUseCase) FindAllExpense(page, size int, startDate, endDate string, user string) ([]entity.Expense, model.Paging, error) {
	return e.repo.List(page, size, startDate, endDate, user)
}

func (e *expenseUseCase) FindExpenseByID(id string) (entity.Expense, error) {
	return e.repo.Get(id)
}

func (e *expenseUseCase) FindExpenseByTransactionType(transactionType string, user string) ([]entity.Expense, error) {
	return e.repo.GetByTransaction(strings.ToUpper(transactionType), user)
}

func NewExpenseUseCase(repo repository.ExpenseRepository) ExpenseUseCase {
	return &expenseUseCase{repo: repo}
}
