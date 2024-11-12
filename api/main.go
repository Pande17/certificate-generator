package main

import (
	"log"
	"pkl/finalProject/certificate-generator/config"
	"pkl/finalProject/certificate-generator/internal/database"
	"pkl/finalProject/certificate-generator/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// intitialize environment variables
	config.InitEnv()

	// connect to database MongoDB
	database.ConnectMongoDB()
	database.CreateCollectionsAndIndexes(database.MongoClient)

	// create a new fiber application instance
	app := fiber.New()

	// setup routes
	routes.RouteSetup(app)

	err := app.Listen(":3000")
	if err != nil {
		log.Println("Error on running fiber, ", err.Error())
	}
}
