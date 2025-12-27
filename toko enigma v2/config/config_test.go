package config

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestNewConfig_Success() {
	// Set environment variables for a successful case
	os.Setenv("DB_HOST", "localhost_test")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("API_PORT", "9090")
	os.Setenv("JWT_SIGNATURE_KEY", "test_secret")

	// Unset them after the test
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("API_PORT")
		os.Unsetenv("JWT_SIGNATURE_KEY")
	}()

	cfg, err := NewConfig()

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cfg)
	assert.Equal(s.T(), "localhost_test", cfg.Host)
	assert.Equal(s.T(), "5433", cfg.Port)
	assert.Equal(s.T(), "test_db", cfg.Database)
	assert.Equal(s.T(), "test_user", cfg.Username)
	assert.Equal(s.T(), "test_password", cfg.Password)
	assert.Equal(s.T(), "9090", cfg.ApiPort)
	assert.Equal(s.T(), []byte("test_secret"), cfg.JwtSignatureKey)
}

func (s *ConfigTestSuite) TestNewConfig_MissingDBPassword() {
	// Set only JWT key, but not DB password
	os.Setenv("JWT_SIGNATURE_KEY", "test_secret")
	defer os.Unsetenv("JWT_SIGNATURE_KEY")

	// Ensure DB_PASSWORD is not set
	os.Unsetenv("DB_PASSWORD")

	cfg, err := NewConfig()

	assert.Error(s.T(), err)
	assert.Nil(s.T(), cfg)
	assert.EqualError(s.T(), err, "environment variable DB_PASSWORD is not set")
}

func (s *ConfigTestSuite) TestNewConfig_MissingJWTSignatureKey() {
	// Set only DB password, but not JWT key
	os.Setenv("DB_PASSWORD", "test_password")
	defer os.Unsetenv("DB_PASSWORD")

	// Ensure JWT_SIGNATURE_KEY is not set
	os.Unsetenv("JWT_SIGNATURE_KEY")

	cfg, err := NewConfig()

	assert.Error(s.T(), err)
	assert.Nil(s.T(), cfg)
	assert.EqualError(s.T(), err, "environment variable JWT_SIGNATURE_KEY is not set")
}

func (s *ConfigTestSuite) TestNewConfig_Defaults() {
	// Set only the required variables
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("JWT_SIGNATURE_KEY", "test_secret")
	defer os.Unsetenv("DB_PASSWORD")
	defer os.Unsetenv("JWT_SIGNATURE_KEY")

	cfg, err := NewConfig()

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cfg)
	// Check default values
	assert.Equal(s.T(), "8080", cfg.ApiPort)
	assert.Equal(s.T(), "Enigma Camp", cfg.ApplicationName)
	assert.Equal(s.T(), jwt.SigningMethodHS256, cfg.JwtSignedMethod)
	assert.Equal(s.T(), time.Hour*1, cfg.AccessTokenLifetime)
}
