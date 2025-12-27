package repository

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/mini-banking/model"
)

type UserRepository interface {
	Create(user model.UserCredential) (model.UserCredential, error)
	List() ([]model.UserCredential, error)
	Get(id uint32) (model.UserCredential, error)
	GetByUsername(username string) (model.UserCredential, error)
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) Create(user model.UserCredential) (model.UserCredential, error) {
	var id uint32

	err := u.db.QueryRow("INSERT INTO mst_user (username, password, role, balance) VALUES  ($1, $2, $3, $4) RETURNING id", user.Username, user.Password, user.Role, user.Balance).Scan(&id)

	if err != nil {
		// fmt.Println(err)
		return model.UserCredential{}, fmt.Errorf("failed to save user")
	}

	user.Id = id

	return user, nil
}

func (u *userRepository) List() ([]model.UserCredential, error) {
	var users []model.UserCredential
	rows, err := u.db.Query("SELECT id, username, role, balance FROM mst_user")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve list user")
	}
	for rows.Next() {
		var user model.UserCredential
		err := rows.Scan(&user.Id, &user.Username, &user.Role, &user.Balance)
		if err != nil {
			return nil, fmt.Errorf("failed to scaning data")
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *userRepository) Get(id uint32) (model.UserCredential, error) {
	var user model.UserCredential

	err := u.db.QueryRow("SELECT id, username, role, balance FROM mst_user WHERE id = $1", id).Scan(&user.Id, &user.Username, &user.Role, &user.Balance)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return model.UserCredential{}, err
		}
		return model.UserCredential{}, fmt.Errorf("failed to get user by ID")
	}
	
	return user, nil
}

func (u *userRepository) GetByUsername(username string) (model.UserCredential, error) {
	var user model.UserCredential

	err := u.db.QueryRow("SELECT id, username, password, role, balance FROM mst_user WHERE username = $1", username).Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.Balance)
	
	if err != nil {
		return model.UserCredential{}, err
	}

	return user, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}
