package config

import (
	"log"

	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Failed loading .env file, using system env")
	}
}
