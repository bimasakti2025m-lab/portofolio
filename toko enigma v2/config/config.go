package config

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type APIConfig struct {
	ApiPort string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     []byte
	JwtSignedMethod     *jwt.SigningMethodHMAC
	AccessTokenLifetime time.Duration
}
type Config struct {
	DBConfig
	APIConfig
	TokenConfig
}



func (c *Config) readConfig() error {
	// Set default values
	c.DBConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.APIConfig = APIConfig{
		ApiPort: "8080",
	}

	c.TokenConfig = TokenConfig{
		ApplicationName:     "Enigma Camp",
		JwtSignatureKey:     []byte(""), // Will be read from env
		JwtSignedMethod:     jwt.SigningMethodHS256,
		AccessTokenLifetime: time.Hour * 1,
	}

	// Read from environment variables and override defaults
	if host := os.Getenv("DB_HOST"); host != "" {
		c.DBConfig.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		c.DBConfig.Port = port
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		c.DBConfig.Database = dbName
	}
	if user := os.Getenv("DB_USER"); user != "" {
		c.DBConfig.Username = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		c.DBConfig.Password = password
	}
	if apiPort := os.Getenv("API_PORT"); apiPort != "" {
		c.APIConfig.ApiPort = apiPort
	}
	if signatureKey := os.Getenv("JWT_SIGNATURE_KEY"); signatureKey != "" {
		c.TokenConfig.JwtSignatureKey = []byte(signatureKey)
	}

	// Validate required secret fields
	if c.DBConfig.Password == "" {
		return fmt.Errorf("environment variable DB_PASSWORD is not set")
	}
	if len(c.TokenConfig.JwtSignatureKey) == 0 {
		return fmt.Errorf("environment variable JWT_SIGNATURE_KEY is not set")
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}
