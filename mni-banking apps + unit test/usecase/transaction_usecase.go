package usecase

import (
	"enigmacamp.com/mini-banking/model"
	"enigmacamp.com/mini-banking/repository"
)

type TransactionUseCase interface {
	CreateTransaction(payload model.Transaction) (model.Transaction, error)
	ListTransaction() ([]model.Transaction, error)
	GetTransactionById(id uint32) (model.Transaction, error)
	GetTransactionByUserId(userId uint32) ([]model.Transaction, error)
	UpdateTransaction(transaction model.Transaction) (model.Transaction, error)
	DeleteTransaction(id uint32) error
}

type transactionUseCase struct {
	repo repository.TransactionRepository
}

func NewTransactionUseCase(repo repository.TransactionRepository) TransactionUseCase {
	return &transactionUseCase{repo: repo}
}

func (t *transactionUseCase) CreateTransaction(payload model.Transaction) (model.Transaction, error) {
	return t.repo.Create(payload)
}

func (t *transactionUseCase) ListTransaction() ([]model.Transaction, error) {
	return t.repo.List()
}

func (t *transactionUseCase) GetTransactionById(id uint32) (model.Transaction, error) {
	return t.repo.Get(id)
}

func (t *transactionUseCase) GetTransactionByUserId(userId uint32) ([]model.Transaction, error) {
	return t.repo.GetByUserId(userId)
}

func (t *transactionUseCase) UpdateTransaction(transaction model.Transaction) (model.Transaction, error) {
	return t.repo.Update(transaction)
}

func (t *transactionUseCase) DeleteTransaction(id uint32) error {
	return t.repo.Delete(id)
}



