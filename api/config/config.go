package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	log.Println("make dir temp err: ", os.Mkdir("temp", 0644))
	log.Println("make dir assets/certificate err: ", os.Mkdir("assets/certificate", 0644))
	if err := godotenv.Load(".env", ".env.local"); err != nil {
		log.Println("Failed loading .env file, using system env")
	}
}
