package main

import (
	"log"
	"pkl/finalProject/certificate-generator/repository/config"
	"pkl/finalProject/certificate-generator/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// intitialize environment variables
	InitEnv()

	// connect to database MongoDB
	config.ConnectMongoDB()

	// create a new fiber application instance
	app := fiber.New()

	// setup routes
	routes.RouteSetup(app)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(
			"Error on running fiber, ", 
			err.Error())
	}


}