// TODO : Membuat repository user terhubung dengan database

package repository

import (
	"E-commerce-Sederhana/model"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	GetUserByUsername(username string) (model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func (ur *userRepository) Create(user *model.User) (*model.User, error) {
	query := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"
	err := ur.db.QueryRow(query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) GetAllUsers() ([]model.User, error) {
	query := "SELECT id, username, email, password, role FROM users"
	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *userRepository) GetUserByUsername(username string) (model.User, error) {
	query := "SELECT id, username, email, password, role FROM users WHERE username = $1"
	var user model.User

	err := ur.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, fmt.Errorf("user with username %s not found", username)
		}
		return model.User{}, err
	}

	return user, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
