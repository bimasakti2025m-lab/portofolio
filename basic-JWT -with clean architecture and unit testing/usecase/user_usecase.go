package usecase

import (
	"basic-JWT/model"
	"basic-JWT/repository"
)

type UserUsecase interface {
	Create(user *model.User) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	GetUserByUsername(username string) (model.User, error)
}

type userUsecase struct {
	userRepository repository.UserRepository
}

func (uu *userUsecase) Create(user *model.User) (*model.User, error) {
	return uu.userRepository.Create(user)
}

func (uu *userUsecase) GetAllUsers() ([]model.User, error) {
	return uu.userRepository.GetAllUsers()
}

func (uu *userUsecase) GetUserByUsername(username string) (model.User, error) {
	return uu.userRepository.GetUserByUsername(username)
}

func NewUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}