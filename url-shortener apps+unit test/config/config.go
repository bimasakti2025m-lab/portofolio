package config

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

type Config struct {
	DBConfig
	APIConfig
	TokenConfig
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     []byte
	JwtSignedMethod     *jwt.SigningMethodHMAC
	AccessTokenLifetime time.Duration
}



func (c *Config) readConfig() error {
	// Read from environment variables
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Set default values
	c.DBConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Driver:   os.Getenv("DB_DRIVER"),
	}
		
	c.APIConfig = APIConfig{
		ApiPort: os.Getenv("API_PORT"),
	}

	tokenLifetime, _ := time.ParseDuration(os.Getenv("TOKEN_ACCESS_TOKEN_LIFETIME"))
	c.TokenConfig = TokenConfig{
		ApplicationName:     os.Getenv("APPLICATION_NAME"),
		JwtSignatureKey:     []byte(os.Getenv("JWT_SIGNATURE_KEY")), // Will be read from env
		JwtSignedMethod:     jwt.SigningMethodHS256,
		AccessTokenLifetime: tokenLifetime,
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
