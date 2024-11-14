package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	log.Println(os.Mkdir("temp", 0644))
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Failed loading .env file, using system env")
	}
}
