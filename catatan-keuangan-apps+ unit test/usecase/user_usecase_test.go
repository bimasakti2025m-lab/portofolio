package usecase

import (
	"errors"
	"testing"

	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/mock/usecase_mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUCSuite struct {
	suite.Suite
	userRepo *usecase_mock.UserUsecaseMock
	userUC   UserUseCase
}

func TestUserUCSuite(t *testing.T) {
	suite.Run(t, new(UserUCSuite))
}

func (u *UserUCSuite) SetupTest() {
	u.userRepo = new(usecase_mock.UserUsecaseMock)
	u.userUC = NewUserUseCase(u.userRepo)
}

func (u *UserUCSuite) TestRegisterNewUser_success() {
	newUser := entity.User{
		Username: "success",
		Password: "password success",
		Role:     "user",
		Expenses: nil,
	}

	u.userRepo.On("GetByUsername", newUser.Username).Return(entity.User{}, errors.New("not found")).Once()
	// Expect Create to be called with hashed password
	u.userRepo.On("Create", mock.MatchedBy(func(user entity.User) bool {
		return user.Username == newUser.Username && user.Role == "user" && user.Password != newUser.Password
	})).Return(entity.User{
		ID:       "uuid-user-test",
		Username: newUser.Username,
		Password: "hashed-password",
		Role:     "user",
	}, nil).Once()

	user, err := u.userUC.RegisterNewUser(newUser)
	u.Equal(user, entity.User{
		ID:       "uuid-user-test",
		Username: newUser.Username,
		Password: "hashed-password",
		Role:     "user",
	})
	u.Nil(err)
}

func (u *UserUCSuite) TestCreate_failed() {
	newUser := entity.User{
		Username: "failed",
		Password: "password failed",
	}

	u.userRepo.On("GetByUsername", newUser.Username).Return(entity.User{Username: newUser.Username}, nil).Once()

	user, err := u.userUC.RegisterNewUser(newUser)
	u.Equal(user, entity.User{})
	u.NotNil(err)
}

func (u *UserUCSuite) TestFindByID_success() {
	existingUser := entity.User{
		ID:       "uuid-user-test",
		Username: "success",
		Password: "password success",
	}

	u.userRepo.On("Get", existingUser.ID).Return(existingUser, nil).Once()

	user, err := u.userUC.FindUserByID(existingUser.ID)
	u.Equal(user, existingUser)
	u.Nil(err)
}

func (u *UserUCSuite) TestFindByID_failed() {
	u.userRepo.On("Get", "uuid-user-failed").Return(entity.User{}, errors.New("not found")).Once()

	user, err := u.userUC.FindUserByID("uuid-user-failed")
	u.Equal(user, entity.User{})
	u.NotNil(err)
}

// func (u *UserUCSuite) TestFindByUsernamePassword_success() {
// 	existingUser := entity.User{
// 		ID:       "uuid-user-test",
// 		Username: "success",
// 		Password: "$2a$10$7a8b9c0d1e2f3g4h5i6j7u8v9w0x1y2z3A4B5C6D7E8F9G0H1I2J3K", // bcrypt hash for "password success"
// 	}

// 	u.userRepo.On("GetByUsername", existingUser.Username).Return(existingUser, nil).Once()

// 	foundUser, err := u.userUC.FindUserByUsernamePassword(existingUser.Username, "password success")
// 	log.Println(foundUser)
// 	u.Nil(err)
// 	// The FindUserByUsernamePassword function returns the user with the hashed password from the repository.
// 	// So, the comparison should be against the existingUser with its hashed password.
// 	// We expect the returned user to be identical to the existingUser we mocked.
// 	u.Equal(existingUser, foundUser)
// }
