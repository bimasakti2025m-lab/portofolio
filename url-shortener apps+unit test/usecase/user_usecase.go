package usecase

import (
	"fmt"

	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterNewUser(payload model.UserCredential) (model.UserCredential, error)
	FindAllUser() ([]model.UserCredential, error)
	FindUserById(id uint32) (model.UserCredential, error)
	FindUserByUsernamePassword(username string, password string) (model.UserCredential, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func (u *userUseCase) RegisterNewUser(payload model.UserCredential) (model.UserCredential, error) {
	_ , err := u.repo.GetByUsernamePassword(payload.Username, payload.Password)
	if err == nil {
		return model.UserCredential{}, fmt.Errorf("user already exists")
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserCredential{}, fmt.Errorf("failed to hash password: %v", err)
	}
	payload.Password = string(bcryptPassword)

	return u.repo.Create(payload)
}

func (u *userUseCase) FindAllUser() ([]model.UserCredential, error) {
	return u.repo.List()
}

func (u *userUseCase) FindUserById(id uint32) (model.UserCredential, error) {
	return u.repo.Get(id)
}

func (u *userUseCase) FindUserByUsernamePassword(username string, password string) (model.UserCredential, error) {
	return u.repo.GetByUsernamePassword(username, password)
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
