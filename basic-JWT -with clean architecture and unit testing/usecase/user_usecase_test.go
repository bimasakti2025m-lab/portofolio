package usecase_test

import (
	"basic-JWT/mock/usecase_mock"
	"basic-JWT/model"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"basic-JWT/usecase"
)

type userUcSuite struct {
	suite.Suite
	userRepo *usecase_mock.UserUseCaseMock
	userUc   usecase.UserUsecase
}

func TestUserUcSuite(t *testing.T) {
	suite.Run(t, new(userUcSuite))
}

func (u *userUcSuite) SetupTest() {
	u.userRepo = new(usecase_mock.UserUseCaseMock)
	u.userUc = usecase.NewUserUsecase(u.userRepo)
}

func (u *userUcSuite) TestCreateUser_Success() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("Create", &user).Return(&user, nil)
	createdUser, err := u.userUc.Create(&user)

	// assert
	u.NoError(err)
	u.Equal(&user, createdUser)
}

func (u *userUcSuite) TestCreateUser_Failed() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("Create", &user).Return(nil, errors.New("error"))

	// assert
	_, err := u.userUc.Create(&user)
	u.EqualError(err, "error")

}

func (u *userUcSuite) TestCreateUser_EmptyUsername() {
	// prepare
	username := ""
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("Create", &user).Return(nil, errors.New("username cannot be empty"))

	// assert
	_, err := u.userUc.Create(&user)
	u.EqualError(err, "username cannot be empty")
}

func (u *userUcSuite) TestCreateUser_EmptyPassword() {
	// prepare
	username := "username"
	password := ""
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("Create", &user).Return(nil, errors.New("password cannot be empty"))

	// assert
	_, err := u.userUc.Create(&user)
	u.EqualError(err, "password cannot be empty")
}

func (u *userUcSuite) TestCreateUser_FailedHashingPassword() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("Create", &user).Return(nil, errors.New("error hashing password"))

	// assert
	_, err := u.userUc.Create(&user)
	u.EqualError(err, "error hashing password")
}

func (u *userUcSuite) TestGetAllUsers_Success() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("GetAllUsers").Return([]model.User{user}, nil)
	users, err := u.userUc.GetAllUsers()

	// assert
	u.NoError(err)
	u.Len(users, 1)
}

func (u *userUcSuite) TestGetAllUsers_Failed() {
	// action
	u.userRepo.On("GetAllUsers").Return([]model.User{}, errors.New("error"))

	// assert
	_, err := u.userUc.GetAllUsers()
	u.EqualError(err, "error")
}

func (u *userUcSuite) TestGetUserByUsername_Success() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("GetUserByUsername", username).Return(user, nil)
	foundUser, err := u.userUc.GetUserByUsername(username)

	// assert
	u.NoError(err)
	u.Equal(user, foundUser)
}

func (u *userUcSuite) TestGetUserByUsername_Failed() {
	// action
	u.userRepo.On("GetUserByUsername", "username").Return(model.User{}, errors.New("error"))

	// assert
	_, err := u.userUc.GetUserByUsername("username")
	u.EqualError(err, "error")
}

func (u *userUcSuite) TestGetUserByUsername_UserNotFound() {
	// action
	u.userRepo.On("GetUserByUsername", "username").Return(model.User{}, errors.New("user with username username not found"))

	// assert
	_, err := u.userUc.GetUserByUsername("username")
	u.EqualError(err, "user with username username not found")
}

func (u *userUcSuite) TestGetUserByUsername_EmptyUsername() {
	// action
	u.userRepo.On("GetUserByUsername", "").Return(model.User{}, errors.New("username cannot be empty"))

	// assert
	_, err := u.userUc.GetUserByUsername("")
	u.EqualError(err, "username cannot be empty")
}

func (u *userUcSuite) TestGetUserByUsername_EmptyPassword() {
	// action
	u.userRepo.On("GetUserByUsername", "username").Return(model.User{}, errors.New("password cannot be empty"))

	// assert
	_, err := u.userUc.GetUserByUsername("username")
	u.EqualError(err, "password cannot be empty")
}

func (u *userUcSuite) TestGetUserByUsername_FailedHashingPassword() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("GetUserByUsername", username).Return(user, errors.New("error hashing password"))

	// assert
	_, err := u.userUc.GetUserByUsername(username)
	u.EqualError(err, "error hashing password")
}

func (u *userUcSuite) TestGetUserByUsername_InvalidPassword() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("GetUserByUsername", username).Return(user, errors.New("invalid password"))

	// assert
	_, err := u.userUc.GetUserByUsername(username)
	u.EqualError(err, "invalid password")
}

func (u *userUcSuite) TestGetUserByUsername_InvalidRole() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("GetUserByUsername", username).Return(user, errors.New("invalid role"))

	// assert
	_, err := u.userUc.GetUserByUsername(username)
	u.EqualError(err, "invalid role")
}

func (u *userUcSuite) TestGetUserByUsername_InvalidID() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("GetUserByUsername", username).Return(user, errors.New("invalid id"))

	// assert
	_, err := u.userUc.GetUserByUsername(username)
	u.EqualError(err, "invalid id")
}

func (u *userUcSuite) TestGetUserByUsername_InvalidUsername() {
	// prepare
	username := "username"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{Username: username, Password: string(hashedPassword), Role: "user"}

	// action
	u.userRepo.On("GetUserByUsername", username).Return(user, errors.New("invalid username"))

	// assert
	_, err := u.userUc.GetUserByUsername(username)
	u.EqualError(err, "invalid username")
}

