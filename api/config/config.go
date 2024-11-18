package config

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func InitEnv() {
	log.Println("make dir temp err: ", os.Mkdir("temp", 0644))
	log.Println("make dir assets/certificate err: ", os.Mkdir("assets/certificate", 0644))
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Failed loading .env file, using system env")
	}

	dirEnt, err := os.ReadDir("/")
	if err != nil {
		log.Println("readdir error: ", err)
	} else {
		readAllDir(dirEnt, 0)
	}
}

func readAllDir(dirEnt []fs.DirEntry, indent int) {
	for _, entry := range dirEnt {
		fmt.Print(strings.Repeat(" ", indent))
		fmt.Print(entry.Name())
		if entry.IsDir() {
			childEnt, err := os.ReadDir("/")
			if err != nil {
				log.Println("readdir error: ", err)
			} else {
				readAllDir(childEnt, indent+2)
			}
		}
		fmt.Println()
	}
}
