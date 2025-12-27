package service

import (
	"testing"
	"time"

	"enigmacamp.com/mini-banking/config"
	"enigmacamp.com/mini-banking/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JwtServiceTestSuite struct {
	suite.Suite
	jwtService JwtService
	tokenCfg   config.TokenConfig
}

func (s *JwtServiceTestSuite) SetupTest() {
	s.tokenCfg = config.TokenConfig{
		ApplicationName:     "test-app",
		JwtSignatureKey:     []byte("my-secret-key"),
		AccessTokenLifetime: 1 * time.Hour,
		JwtSignedMethod:     jwt.SigningMethodHS256,
	}
	s.jwtService = NewJwtService(s.tokenCfg)
}

func TestJwtServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JwtServiceTestSuite))
}

func (s *JwtServiceTestSuite) TestCreateToken_Success() {
	user := model.UserCredential{Id: 1, Role: "admin"}

	token, err := s.jwtService.CreateToken(user)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), token)

	// Verifikasi token yang baru dibuat untuk memastikan isinya benar
	claims, err := s.jwtService.VerifyToken(token)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Id, claims.UserId)
	assert.Equal(s.T(), user.Role, claims.Role)
	assert.Equal(s.T(), s.tokenCfg.ApplicationName, claims.Issuer)
}

func (s *JwtServiceTestSuite) TestVerifyToken_Success() {
	user := model.UserCredential{Id: 5, Role: "user"}
	token, _ := s.jwtService.CreateToken(user)

	claims, err := s.jwtService.VerifyToken(token)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), claims)
	assert.Equal(s.T(), user.Id, claims.UserId)
	assert.Equal(s.T(), user.Role, claims.Role)
}

func (s *JwtServiceTestSuite) TestVerifyToken_Expired() {
	// Buat service baru dengan token yang langsung kedaluwarsa
	expiredTokenCfg := config.TokenConfig{
		ApplicationName:     "test-app",
		JwtSignatureKey:     []byte("my-secret-key"),
		AccessTokenLifetime: -1 * time.Minute, // Waktu di masa lalu
		JwtSignedMethod:     jwt.SigningMethodHS256,
	}
	expiredJwtService := NewJwtService(expiredTokenCfg)

	user := model.UserCredential{Id: 1, Role: "admin"}
	token, _ := expiredJwtService.CreateToken(user)

	claims, err := s.jwtService.VerifyToken(token)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), claims)
	assert.ErrorIs(s.T(), err, jwt.ErrTokenExpired)
}

func (s *JwtServiceTestSuite) TestVerifyToken_InvalidSignature() {
	// Buat service dengan secret key yang berbeda
	anotherJwtService := NewJwtService(config.TokenConfig{
		JwtSignatureKey: []byte("different-secret-key"),
		JwtSignedMethod: jwt.SigningMethodHS256,
	})

	user := model.UserCredential{Id: 1, Role: "admin"}
	// Buat token dengan service lain
	token, _ := anotherJwtService.CreateToken(user)

	// Verifikasi dengan service utama
	claims, err := s.jwtService.VerifyToken(token)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), claims)
	assert.ErrorIs(s.T(), err, jwt.ErrSignatureInvalid)
}

func (s *JwtServiceTestSuite) TestVerifyToken_InvalidFormat() {
	claims, err := s.jwtService.VerifyToken("this.is.not.a.valid.token")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), claims)
	assert.ErrorIs(s.T(), err, jwt.ErrTokenMalformed)
}
