package main

import (
	"certificate-generator/api/config"
	"certificate-generator/api/internal/database"
	"certificate-generator/api/routes"
	"log"

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
