// TODO :
// 1. Mendeklarasikan struct DBConfig, APIConfig, dan TokenConfig
// 2. Mendeklarasikan struct Config
// 3. Mendeklarasikan method readCOnfig
// 4. Mendeklarasikan konstruktor NewConfig

package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Driver   string
}

type APIConfig struct {
	Port string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     []byte
	JwtSignedMethod     *jwt.SigningMethodHMAC
	AccessTokenLifetime time.Duration
}

type MidtransConfig struct {
	ServerKey string `env:"MIDTRANS_SERVER_KEY"`
	Env       string `env:"MIDTRANS_ENV"`
}

type Config struct {
	DB       DBConfig
	API      APIConfig
	Token    TokenConfig
	Midtrans MidtransConfig
}

func (c *Config) readConfig() error {
	// Secara eksplisit memuat file .env dari root proyek
	// Sesuaikan path jika struktur direktori Anda berbeda
	err := godotenv.Load() // godotenv.Load() akan mencari .env di direktori saat ini
	if err != nil {
		// Jangan gunakan Fatalf agar bisa ditangani, tapi beri log yang jelas
		log.Printf("Warning: Could not load .env file. Reading from environment variables. Error: %v", err)
	}

	c.DB.Host = os.Getenv("DB_HOST")
	c.DB.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	c.DB.Username = os.Getenv("DB_USERNAME")
	c.DB.Password = os.Getenv("DB_PASSWORD")
	c.DB.Database = os.Getenv("DB_DATABASE")
	c.DB.Driver = os.Getenv("DB_DRIVER")

	c.API.Port = os.Getenv("API_PORT")

	c.Token.ApplicationName = os.Getenv("TOKEN_APPLICATION_NAME")
	c.Token.JwtSignatureKey = []byte(os.Getenv("TOKEN_JWT_SIGNATURE_KEY"))
	c.Token.JwtSignedMethod = jwt.SigningMethodHS256
	c.Token.AccessTokenLifetime, _ = time.ParseDuration(os.Getenv("TOKEN_ACCESS_TOKEN_LIFETIME"))

	c.Midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	c.Midtrans.Env = os.Getenv("MIDTRANS_ENV")

	// Validasi variabel penting
	if c.Midtrans.ServerKey == "" {
		log.Fatal("FATAL: MIDTRANS_SERVER_KEY environment variable is not set.")
	}

	return nil
}

func NewConfig() (*Config, error) {
	var c Config
	err := c.readConfig()
	if err != nil {
		return nil, err
	}
	return &c, nil

}
