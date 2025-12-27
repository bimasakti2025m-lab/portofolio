package config

const (
	SelectExpenseList              = `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE user_id = $3 ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	SelectExpenseListFull          = `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE date BETWEEN $3 AND $4 AND user_id = $5 ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	SelectExpenseByID              = `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE id = $1`
	SelectExpenseByTransactionType = `SELECT id, date, amount, transaction_type, balance, description, created_at, updated_at FROM expenses WHERE transaction_type=$1 AND user_id = $2 ORDER BY created_at DESC`
	SelectCountExpense             = `SELECT COUNT(*) FROM expenses WHERE user_id=$1`
	InsertExpenses                 = `INSERT INTO expenses (date, amount, transaction_type, balance, description, user_id, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, balance, created_at`
)
