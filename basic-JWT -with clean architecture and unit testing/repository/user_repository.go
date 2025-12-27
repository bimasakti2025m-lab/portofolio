// TODO : Membuat repository user terhubung dengan database

package repository

import (
	"basic-JWT/model"
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
	err := ur.db.QueryRow("INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id", user.Username, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	rows, err := ur.db.Query("SELECT id, username, password, role FROM users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur *userRepository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	row := ur.db.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role); err != nil {
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
