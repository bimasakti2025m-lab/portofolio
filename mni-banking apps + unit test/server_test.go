package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) TestNewServer_Success() {
	// Atur environment variables yang dibutuhkan untuk tes
	// Ini akan menimpa nilai dari .env jika ada
	os.Setenv("DB_PASSWORD", "123qweasd") // Ganti dengan password DB test Anda jika berbeda
	os.Setenv("JWT_SIGNATURE_KEY", "test-secret")

	// Pastikan NewServer() tidak panic
	assert.NotPanics(s.T(), func() {
		server := NewServer()
		// Verifikasi bahwa server dan dependensi utamanya berhasil dibuat
		assert.NotNil(s.T(), server)
		assert.NotNil(s.T(), server.engine)
		assert.NotNil(s.T(), server.userUC)
	})
}

func (s *ServerTestSuite) TestNewServer_Fail() {
	// Atur environment variables yang dibutuhkan untuk tes
	// Ini akan menimpa nilai dari .env jika ada
	os.Setenv("DB_PASSWORD", "123qweasd") // Ganti dengan password DB test Anda jika berbeda
	os.Setenv("JWT_SIGNATURE_KEY", "test-secret")

	// Pastikan NewServer() tidak panic
	assert.NotPanics(s.T(), func() {
		server := NewServer()
		// Verifikasi bahwa server dan dependensi utamanya berhasil dibuat
		assert.NotNil(s.T(), server)
		assert.NotNil(s.T(), server.engine)
		assert.NotNil(s.T(), server.userUC)
	})
}