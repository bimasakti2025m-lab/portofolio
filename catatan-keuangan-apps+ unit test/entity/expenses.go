package entity

import (
	"time"
)

type Expense struct {
	ID              string    `json:"id"`
	Date            time.Time `json:"date"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transactionType"`
	Balance         float64   `json:"balance,omitempty"`
	Description     string    `json:"description"`
	UserId          string    `json:"userId,omitempty"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (e Expense) IsTransactionTypeValid() bool {
	return e.TransactionType == "CREDIT" || e.TransactionType == "DEBIT"
}

func (e Expense) IsRequiredFields() bool {
	return e.Amount > 0 || e.TransactionType != "" || e.Description != ""
}
