package usecase

import (
	"encoding/base64"

	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/entity/dto"
	"enigmacamp.com/livecode-catatan-keuangan/shared/service"
)

type AuthUseCase interface {
	Login(payload dto.AuthRequestDto) (dto.AuthResponseDto, error)
	Register(payload dto.AuthRequestDto) (entity.User, error)
}

type authUseCase struct {
	uc         UserUseCase
	jwtService service.JwtService
}

func (a *authUseCase) Login(payload dto.AuthRequestDto) (dto.AuthResponseDto, error) {
	user, err := a.uc.FindUserByUsernamePassword(payload.Username, payload.Password)
	if err != nil {
		return dto.AuthResponseDto{}, err
	}
	// token, err := a.jwtService.CreateToken(user)
	// if err != nil {
	// 	return dto.AuthResponseDto{}, err
	// }

	// TODO: Create base64 token
	token := base64.StdEncoding.EncodeToString([]byte(user.Role + ":" + user.ID))
	response := dto.AuthResponseDto{
		Token: token,
	}
	return response, nil
}

func (a *authUseCase) Register(payload dto.AuthRequestDto) (entity.User, error) {
	return a.uc.RegisterNewUser(entity.User{Username: payload.Username, Password: payload.Password})
}

func NewAuthUseCase(uc UserUseCase, jwtService service.JwtService) AuthUseCase {
	return &authUseCase{uc: uc, jwtService: jwtService}
}
