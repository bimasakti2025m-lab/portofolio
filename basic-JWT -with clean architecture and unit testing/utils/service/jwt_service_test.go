package service

import (
	"basic-JWT/config"
	"basic-JWT/model"
	"testing"

	"github.com/stretchr/testify/suite"
)

type JWTServiceSuite struct {
	suite.Suite
	jwtSvc JWTservice
}

func TestJWTServiceSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceSuite))
}

func (j *JWTServiceSuite) SetupTest() {
	// Use a dummy config for testing
	cfg := config.TokenConfig{
		JwtSignatureKey: []byte("secret"),
	}
	j.jwtSvc = NewJWTService(cfg)
}

func (j *JWTServiceSuite) TestCreateToken_Success() {
	// Action
	token := j.jwtSvc.CreateToken(model.User{ID: 1, Username: "username", Password: "password", Role: "user"})

	// Assert
	j.NotEmpty(token)
}

func (j *JWTServiceSuite) TestVerifyToken_Success() {
	user := model.User{ID: 1, Username: "username", Role: "user"}
	token := j.jwtSvc.CreateToken(user)

	// Action
	claims, err := j.jwtSvc.VerifyToken(token)

	// Assert
	j.NoError(err)
	j.Equal(user.ID, claims.UserId)
	j.Equal(user.Role, claims.Role)
}

func (j *JWTServiceSuite) TestVerifyToken_InvalidToken() {
	// Action
	claims, err := j.jwtSvc.VerifyToken("invalid.token.string")

	// Assert
	j.Error(err)
	j.Nil(claims)
}
