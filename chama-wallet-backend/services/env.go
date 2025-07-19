package services

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found or failed to load")
	}
}
