package modelutils

import "github.com/golang-jwt/jwt/v5"

type JwtPayloadClaims struct {
	jwt.RegisteredClaims
	UserId string
	Role string
}