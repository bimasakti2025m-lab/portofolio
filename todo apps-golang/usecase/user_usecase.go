package usecase

import (
	"enigmacamp.com/golang-jwt/model"
	"enigmacamp.com/golang-jwt/repository"
)

type UserUseCase interface {
	RegisterNewUser(user model.UserCredential) (model.UserCredential, error)
	FindAllUser() ([]model.UserCredential, error)
	FindUserById(id uint32) (model.UserCredential, error)
	FindUserByUsername(username string) (model.UserCredential, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func (u *userUseCase) RegisterNewUser(user model.UserCredential) (model.UserCredential, error) {
	return u.repo.Create(user)
}

func (u *userUseCase) FindAllUser() ([]model.UserCredential, error) {
	return u.repo.List()
}

func (u *userUseCase) FindUserById(id uint32) (model.UserCredential, error) {
	return u.repo.Get(id)
}

func (u *userUseCase) FindUserByUsername(username string) (model.UserCredential, error) {
	return u.repo.GetByUsername(username)
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
