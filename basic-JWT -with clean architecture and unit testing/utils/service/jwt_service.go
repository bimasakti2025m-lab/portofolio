// TODO :
// 1. Mendeklarasikan interface JWTservice
// 2. Mendeklarasikan struct jwtService
// 3. Mendeklarasikan method createToken dan verifyToken
// 4. Mendeklarasikan konstruktor NewJWTService

package service

import (
	"basic-JWT/config"
	"basic-JWT/model"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
	modelutils "basic-JWT/utils/model_utils"
)

type JWTservice interface {
	CreateToken(user model.User) string
	VerifyToken(tokenString string) (*modelutils.JwtPayloadClaims, error)
}

type jwtService struct {
	tokenConfig config.TokenConfig
}

func NewJWTService(tokenConfig config.TokenConfig) JWTservice {
	return &jwtService{
		tokenConfig: tokenConfig,
	}
}

func (j *jwtService) CreateToken(user model.User) string {
	tokenKey := j.tokenConfig.JwtSignatureKey

	claims := modelutils.JwtPayloadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Enigma Camp",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		UserId: user.ID,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func (j *jwtService) VerifyToken(tokenString string) (*modelutils.JwtPayloadClaims, error) {
	tokenKey := j.tokenConfig.JwtSignatureKey
	claims := &modelutils.JwtPayloadClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
