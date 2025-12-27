package repository

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/mini-banking/model"
)

type TransactionRepository interface {
	Create(transaction model.Transaction) (model.Transaction, error)
	List() ([]model.Transaction, error)
	Get(id uint32) (model.Transaction, error)
	GetByUserId(userId uint32) ([]model.Transaction, error)
	Update(transaction model.Transaction) (model.Transaction, error)
	Delete(id uint32) error
}

type transactionRepository struct {
	db *sql.DB
}

// Mendeklarasikan method yang mengimplmentasikan interface TransactionRepository
func (t *transactionRepository) Create(transaction model.Transaction) (model.Transaction, error) {
	var id uint
	err := t.db.QueryRow("INSERT INTO mst_transaction (from_user_id, to_user_id, amount, type, status) VALUES ($1, $2, $3, $4, $5) RETURNING id", transaction.FromUserID, transaction.ToUserID, transaction.Amount, transaction.Type, transaction.Status).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return model.Transaction{}, err
	}
	transaction.ID = id
	return transaction, nil
	}

func (t *transactionRepository) List() ([]model.Transaction, error) {
	var transactions []model.Transaction
	rows, err := t.db.Query("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.FromUserID, &transaction.ToUserID, &transaction.Amount, &transaction.Type)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (t *transactionRepository) Get(id uint32) (model.Transaction, error) {
	var transaction model.Transaction
	err := t.db.QueryRow("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction WHERE id = $1", id).Scan(&transaction.ID, &transaction.FromUserID, &transaction.ToUserID, &transaction.Amount, &transaction.Type)
	if err != nil {
		return model.Transaction{}, err
	}
	return transaction, nil

}

func (t *transactionRepository) GetByUserId(userId uint32) ([]model.Transaction, error) {
	var transactions []model.Transaction
	rows, err := t.db.Query("SELECT id, from_user_id, to_user_id, amount, type FROM mst_transaction WHERE from_user_id = $1 OR to_user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.FromUserID, &transaction.ToUserID, &transaction.Amount, &transaction.Type)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (t *transactionRepository) Update(transaction model.Transaction) (model.Transaction, error) {
	_, err := t.db.Exec("UPDATE mst_transaction SET from_user_id = $1, to_user_id = $2, amount = $3, type = $4 WHERE id = $5", transaction.FromUserID, transaction.ToUserID, transaction.Amount, transaction.Type, transaction.ID)
	if err != nil {
		return model.Transaction{}, err
	}
	return transaction, nil
}

func (t *transactionRepository) Delete(id uint32) error {
	_, err := t.db.Exec("DELETE FROM mst_transaction WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// Mendeklarasikan konstruktor
func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}
