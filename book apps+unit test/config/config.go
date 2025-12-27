// TODO :
// 1. Mendeklarasikan nama package config
// 2. Mendeklarasikan struct DBConfig dan APIConfig
// 3. Mendeklarasikan function baru bernama readConfig
// 4. Mendeklarasikan function baru bernama newConfig

package config

import (
	"errors"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type APIConfig struct {
	Host string
	Port int
}

type Config struct {
	DBConfig
	APIConfig
}

func (c *Config) readConfig() error {
	c.DBConfig = DBConfig{
		Host:     getEnv("DB_HOST", ""),
		Port:     getEnvInt("DB_PORT", 0),
		Username: getEnv("DB_USERNAME", ""),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_DATABASE", ""),
	}
	c.APIConfig = APIConfig{
		Port: getEnvInt("API_PORT", 0),
	}
	if c.DBConfig.Host == "" || c.DBConfig.Port == 0 || c.DBConfig.Username == "" || c.DBConfig.Password == "" || c.DBConfig.Database == "" {
		return errors.New("must be filled")
	}
	return nil
}

func NewConfig() (*Config, error) {
	c := &Config{}
	if err := c.readConfig(); err != nil {
		return nil, err
	}
	return c, nil
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}