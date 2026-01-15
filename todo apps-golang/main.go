package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Golang JWT")
	// Muat file .env jika ada
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file")
		}
	} else {
		log.Println("No .env file found. Using system environment variables.")
	}
	NewServer().Run()
}
