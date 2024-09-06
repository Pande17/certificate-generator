package routes

import (
	"pkl/finalProject/certificate-generator/controller"
	"pkl/finalProject/certificate-generator/middleware"

	"github.com/gofiber/fiber/v2"
)

func RouteSetup(r *fiber.App) {
	// define a route group for api routes
	api := r.Group("")

	api.Post("/signup", controller.SignUp) // routes for signup admin
	api.Post("/login", controller.Login)	// routes for login admin
	api.Get("/validate", middleware.RequireAuth, controller.Validate)


	// define a route group for template routes
	// template := r.Group("")

	// template.Get("/") // routes for homepage html

}