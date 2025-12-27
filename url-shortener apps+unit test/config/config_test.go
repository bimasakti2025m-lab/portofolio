package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	// Simpan env vars asli untuk dipulihkan setelah tes
	originalEnv map[string]string
}

func (s *ConfigTestSuite) SetupTest() {
	// Simpan env vars yang ada sebelum diubah oleh tes
	s.originalEnv = make(map[string]string)
	envKeys := []string{
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
		"API_PORT", "JWT_SIGNATURE_KEY",
	}
	for _, key := range envKeys {
		if val, ok := os.LookupEnv(key); ok {
			s.originalEnv[key] = val
		}
		os.Unsetenv(key) // Bersihkan env var sebelum setiap tes
	}
}

func (s *ConfigTestSuite) TearDownTest() {
	// Pulihkan env vars asli setelah setiap tes selesai
	for key, val := range s.originalEnv {
		os.Setenv(key, val)
	}
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestNewConfig_SuccessWithOverrides() {
	// Atur environment variables untuk skenario sukses
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "1234")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("API_PORT", "9090")
	os.Setenv("JWT_SIGNATURE_KEY", "testsecret")

	cfg, err := NewConfig()

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cfg)

	// Verifikasi bahwa nilai dari env var telah digunakan
	assert.Equal(s.T(), "testhost", cfg.DBConfig.Host)
	assert.Equal(s.T(), "1234", cfg.DBConfig.Port)
	assert.Equal(s.T(), "testdb", cfg.DBConfig.Database)
	assert.Equal(s.T(), "testuser", cfg.DBConfig.Username)
	assert.Equal(s.T(), "testpass", cfg.DBConfig.Password)
	assert.Equal(s.T(), "9090", cfg.APIConfig.ApiPort)
	assert.Equal(s.T(), []byte("testsecret"), cfg.TokenConfig.JwtSignatureKey)
	assert.Equal(s.T(), time.Hour*1, cfg.TokenConfig.AccessTokenLifetime) // Verifikasi nilai default
}

func (s *ConfigTestSuite) TestNewConfig_SuccessWithDefaults() {
	// Hanya atur variabel yang wajib diisi
	os.Setenv("DB_PASSWORD", "defaultpass")
	os.Setenv("JWT_SIGNATURE_KEY", "defaultsecret")

	cfg, err := NewConfig()

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cfg)

	// Verifikasi bahwa nilai default digunakan
	assert.Equal(s.T(), "localhost", cfg.DBConfig.Host)
	assert.Equal(s.T(), "5432", cfg.DBConfig.Port)
	assert.Equal(s.T(), "8080", cfg.APIConfig.ApiPort)
}

func (s *ConfigTestSuite) TestNewConfig_Fail_MissingDbPassword() {
	os.Setenv("JWT_SIGNATURE_KEY", "testsecret") // JWT key ada
	// DB_PASSWORD sengaja tidak di-set

	_, err := NewConfig()
	assert.Error(s.T(), err)
	assert.EqualError(s.T(), err, "environment variable DB_PASSWORD is not set")
}

func (s *ConfigTestSuite) TestNewConfig_Fail_MissingJwtKey() {
	os.Setenv("DB_PASSWORD", "testpass") // DB password ada
	// JWT_SIGNATURE_KEY sengaja tidak di-set

	_, err := NewConfig()
	assert.Error(s.T(), err)
	assert.EqualError(s.T(), err, "environment variable JWT_SIGNATURE_KEY is not set")
}
