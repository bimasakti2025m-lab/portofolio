package usecase

import (
	"fmt"
	"time"

	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterNewUser(payload entity.User) (entity.User, error)
	FindUserByID(id string) (entity.User, error)
	FindUserByUsernamePassword(username, password string) (entity.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func (u *userUseCase) RegisterNewUser(payload entity.User) (entity.User, error) {
	userExist, _ := u.repo.GetByUsername(payload.Username)
	if userExist.Username == payload.Username {
		return entity.User{}, fmt.Errorf("user with username: %s already exists", payload.Username)
	}
	payload.Role = "user"
	payload.UpdatedAt = time.Now()	

	// TODO: hash password bycrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, err
	}
	payload.Password = string(hash)

	return u.repo.Create(payload)
}

func (u *userUseCase) FindUserByID(id string) (entity.User, error) {
	return u.repo.Get(id)
}

func (u *userUseCase) FindUserByUsernamePassword(username, password string) (entity.User, error) {
	userExist, err := u.repo.GetByUsername(username)
	if err != nil {
		return entity.User{}, fmt.Errorf("user doesn't exists")
	}

	// Cara Cupu!
	// if userExist.Password != password {
	// 	return entity.User{}, fmt.Errorf("password doesn't match")
	// }

	// TODO: compare password bycrypt
	err = bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(password))
	if err != nil {
		return entity.User{}, fmt.Errorf("password doesn't match")
	}

	return userExist, nil
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
