package usecase

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/mini-banking/model"
	"enigmacamp.com/mini-banking/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterNewUser(payload model.UserCredential) (model.UserCredential, error)
	FindAllUser() ([]model.UserCredential, error)
	FindUserById(id uint32) (model.UserCredential, error)
	FindUserByUsernamePassword(username string, password string) (model.UserCredential, error) // Kept for interface compatibility, but implementation changes
}

type userUseCase struct {
	repo repository.UserRepository
}

func (u *userUseCase) RegisterNewUser(payload model.UserCredential) (model.UserCredential, error) {
	return u.repo.Create(payload)
}

func (u *userUseCase) FindAllUser() ([]model.UserCredential, error) {
	return u.repo.List()
}

func (u *userUseCase) FindUserById(id uint32) (model.UserCredential, error) {
	return u.repo.Get(id)
}

func (u *userUseCase) FindUserByUsernamePassword(username string, password string) (model.UserCredential, error) {
	user, err := u.repo.GetByUsername(username)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return model.UserCredential{}, fmt.Errorf("invalid username or password")
		}
		return model.UserCredential{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.UserCredential{}, fmt.Errorf("invalid username or password")
	}

	return user, nil
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
