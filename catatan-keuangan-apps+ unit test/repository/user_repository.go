package repository

import (
	"database/sql"
	"time"

	"enigmacamp.com/livecode-catatan-keuangan/entity"
)

type UserRepository interface {
	Create(payload entity.User) (entity.User, error)
	Get(id string) (entity.User, error)
	GetByUsername(username string) (entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) Create(payload entity.User) (entity.User, error) {
	err := u.db.QueryRow(`INSERT INTO users (username, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`, payload.Username, payload.Password, payload.Role, payload.CreatedAt, payload.UpdatedAt).Scan(&payload.ID, &payload.CreatedAt, &payload.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	return payload, nil
}

func (u *userRepository) Get(id string) (entity.User, error) {
	var user entity.User
	if err := u.db.QueryRow(`SELECT id, username, role, created_at, updated_at FROM users WHERE id=$1`, id).Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return entity.User{}, err
	}
	rows, err := u.db.Query(`SELECT u.id, e.id, e.date, e.amount, e.transaction_type,e.balance, e.description, e.created_at, e.updated_at
FROM users u JOIN expenses e on u.id = e.user_id WHERE u.id = $1`, user.ID)
	if err != nil {
		return entity.User{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var expense entity.Expense
		if err := rows.Scan(&user.ID, &expense.ID, &expense.Date, &expense.Amount, &expense.TransactionType, &expense.Balance, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt); err != nil {
			return entity.User{}, err
		}

		user.Expenses = append(user.Expenses, expense)
	}
	return user, nil
}

func (u *userRepository) GetByUsername(username string) (entity.User, error) {
	var user entity.User
	if err := u.db.QueryRow("SELECT id, username, password, role, created_at, updated_at FROM users WHERE username= $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}
