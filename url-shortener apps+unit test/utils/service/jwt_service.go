package service

import (
	"time"

	"enigmacamp.com/url-shortener/config"
	"enigmacamp.com/url-shortener/model"
	modelutils "enigmacamp.com/url-shortener/utils/model_utils"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(user model.UserCredential) (string, error)
	VerifyToken(token string) (*modelutils.JwtPayloadClaims, error)
}

type jwtService struct {
	tokenConfig config.TokenConfig
}

func (j *jwtService) CreateToken(user model.UserCredential) (string, error) {
	tokenKey := j.tokenConfig.JwtSignatureKey

	claims := modelutils.JwtPayloadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.tokenConfig.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenConfig.AccessTokenLifetime)),
		},
		UserId: user.Id,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(j.tokenConfig.JwtSignedMethod, claims)
	tokenString, err := token.SignedString(tokenKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwtService) VerifyToken(token string) (*modelutils.JwtPayloadClaims, error) {
	tokenKey := j.tokenConfig.JwtSignatureKey
	claims := &modelutils.JwtPayloadClaims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		// Memastikan algoritma penandatanganan token sesuai dengan yang diharapkan (HS256)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return tokenKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}
	return claims, nil
}

func NewJwtService(tokenConfig config.TokenConfig) JwtService {
	return &jwtService{
		tokenConfig: tokenConfig,
	}
}
