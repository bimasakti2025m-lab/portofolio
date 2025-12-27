package repository

import (
	"database/sql"
	"enigmacamp.com/livecode-catatan-keuangan/config"
	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/shared/model"
	"log"
	"math"
)

type ExpenseRepository interface {
	Create(payload entity.Expense) (entity.Expense, error)
	List(page, size int, startDate, endDate string, user string) ([]entity.Expense, model.Paging, error)
	Get(id string) (entity.Expense, error)
	GetByTransaction(transactionType string, user string) ([]entity.Expense, error)
	GetBalance(user string) (float64, error)
}

type expenseRepository struct {
	db *sql.DB
}

func (e *expenseRepository) GetBalance(user string) (float64, error) {
	var balance float64
	if err := e.db.QueryRow(`SELECT balance FROM expenses WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`, user).Scan(&balance); err != nil {
		log.Printf("ExpenseRepository.GetBalance: %v \n", err.Error())
		return 0, err
	}
	return balance, nil
}

func (e *expenseRepository) Create(payload entity.Expense) (entity.Expense, error) {
	err := e.db.QueryRow(config.InsertExpenses, payload.Date, payload.Amount, payload.TransactionType, payload.Balance, payload.Description, payload.UserId, payload.UpdatedAt).Scan(&payload.ID, &payload.Balance, &payload.CreatedAt)
	if err != nil {
		log.Printf("ExpenseRepository.Create: %v \n", err.Error())
		return entity.Expense{}, err
	}
	return payload, nil
}

func (e *expenseRepository) List(page, size int, startDate, endDate string, user string) ([]entity.Expense, model.Paging, error) {
	var expenses []entity.Expense
	offset := (page - 1) * size
	var rows *sql.Rows
	var err error
	if startDate != "" && endDate != "" {
		rows, err = e.db.Query(config.SelectExpenseListFull, size, offset, startDate, endDate, user)
	} else {
		rows, err = e.db.Query(config.SelectExpenseList, size, offset, user)
	}

	if err != nil {
		log.Printf("ExpenseRepository.List: %v \n", err.Error())
		return nil, model.Paging{}, err
	}
	for rows.Next() {
		var expense entity.Expense
		if err := rows.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt); err != nil {
			log.Printf("ExpenseRepository.List.Rows.Next(): %v \n", err.Error())
			return nil, model.Paging{}, err
		}
		expenses = append(expenses, expense)
	}
	totalRows := 0
	if err := e.db.QueryRow(config.SelectCountExpense, user).Scan(&totalRows); err != nil {
		log.Printf("ExpenseRepository.List.Count: %v \n", err.Error())
		return nil, model.Paging{}, err
	}
	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(size))),
	}
	return expenses, paging, nil
}

func (e *expenseRepository) Get(id string) (entity.Expense, error) {
	var expense entity.Expense
	if err := e.db.QueryRow(config.SelectExpenseByID, id).Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt); err != nil {
		log.Printf("ExpenseRepository.Get: %v \n", err.Error())
		return entity.Expense{}, err
	}
	return expense, nil
}

func (e *expenseRepository) GetByTransaction(transactionType string, user string) ([]entity.Expense, error) {
	var expenses []entity.Expense
	rows, err := e.db.Query(config.SelectExpenseByTransactionType, transactionType, user)
	if err != nil {
		log.Printf("ExpenseRepository.GetByTransaction: %v \n", err.Error())
		return nil, err
	}
	for rows.Next() {
		var expense entity.Expense
		if err := rows.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt); err != nil {
			log.Printf("ExpenseRepository.GetByTransaction.Rows.Next(): %v \n", err.Error())
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}
