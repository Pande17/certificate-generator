package routes

import (
	"pkl/finalProject/certificate-generator/controller"
	"pkl/finalProject/certificate-generator/middleware"

	"github.com/gofiber/fiber/v2"
)

// function for template routes
func TemplateRoute(r *fiber.App) {
	accounts := r.Group("/accounts") // define routes group for authentication

	accounts.Get("/login", controller.TEMPlate)  // route for template login page
	accounts.Get("/logout", controller.TEMPlate) // route for template logout

	protected := r.Use(middleware.ValidateCookie) // Define protected routes

	protected.Get("/", controller.TEMPlate)                   // route for template homepage (display all certificate)
	protected.Get("/create-certificate", controller.TEMPlate) // route for template create certificate
	protected.Get("/add-competence", controller.TEMPlate)     // route for template add competence
}
