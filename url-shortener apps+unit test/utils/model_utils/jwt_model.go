package modelutils

import "github.com/golang-jwt/jwt/v5"

type JwtPayloadClaims struct {
	jwt.RegisteredClaims
	UserId uint32
	Role string
}