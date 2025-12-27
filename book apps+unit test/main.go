package main

import (
	"log"

	"github.com/joho/godotenv"
	// Import package konfigurasi Anda
)

func main() {
	// 1. Muat file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Pastikan file .env ada dan isinya benar:", err)
	}

	// 2. Lanjutkan proses membaca konfigurasi dan inisialisasi server

	NewServer().Run()
}
