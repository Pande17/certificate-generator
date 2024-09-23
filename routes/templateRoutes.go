package routes

import (
	"pkl/finalProject/certificate-generator/middleware"

	"github.com/gofiber/fiber/v2"
)

// function for template routes
func TemplateRoute(r *fiber.App) {
	accounts := r.Group("/accounts") // define routes group for authentication

	accounts.Get("/login", nil)  // route for template login page
	accounts.Get("/logout", nil) // route for template logout

	protected := r.Use(middleware.ValidateCookie) // Define protected routes

	protected.Get("/", nil)                   // route for template homepage (display all certificate)
	protected.Get("/create-certificate", nil) // route for template create certificate
	protected.Get("/add-competence", nil)     // route for template add competence
}
