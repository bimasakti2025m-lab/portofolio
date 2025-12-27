package usecase_mock

import (
	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/entity/dto"
	"github.com/stretchr/testify/mock"
)

type AuthUsecaseMock struct {
	mock.Mock
}

func (a *AuthUsecaseMock) Login(payload dto.AuthRequestDto) (dto.AuthResponseDto, error) {
	args := a.Called(payload)
	return args.Get(0).(dto.AuthResponseDto), args.Error(1)
}

func (a *AuthUsecaseMock) Register(payload dto.AuthRequestDto) (entity.User, error) {
	args := a.Called(payload)
	return args.Get(0).(entity.User), args.Error(1)
}
