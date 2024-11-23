package main

import (
	"certificate-generator/config"
	"certificate-generator/database"
	"certificate-generator/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitEnv()

	database.ConnectMongoDB()
	database.CreateCollectionsAndIndexes(database.MongoClient)

	app := fiber.New()
	routes.RouteSetup(app)

	log.Fatal(app.Listen(":3000"))
}
