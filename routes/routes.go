package routes

import "github.com/gofiber/fiber/v2"

func Routes(r *fiber.App) {
	// define a route group for organizing the routes
	certGroup := r.Group("")

	// routes for login admin
	certGroup.Post("/login")	
}