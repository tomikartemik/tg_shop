package main

import (
	"github.com/joho/godotenv"
	"log"
	"tg_shop/internal/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	_, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Application started successfully")
}
