package routes

import (
	"pkl/finalProject/certificate-generator/controller"
	"pkl/finalProject/certificate-generator/middleware"

	"github.com/gofiber/fiber/v2"
)

func RouteSetup(r *fiber.App) {
	r.Get("/", func (c *fiber.Ctx) error  {
		return c.JSON(fiber.Map{
			"message": "test",
		})
	})

	// Define a group routes for API
	api := r.Group("/api")

	// Define routes for authentication
	api.Post("/signup", controller.SignUp)  // Route for signing up admin
	api.Post("/login", controller.Login)    // Route for admin login
	api.Get("/validate", middleware.ValidateCookie, controller.Validate) // Route to check cookie from admin

	// Define protected routes
	// Setiap request ke path dengan group "protected" selalu cek cookie
	protected := api.Use(middleware.ValidateCookie)

	// Define routes for authentication
	protected.Post("/logout", controller.Logout) // Route to logout from account

	// Define routes for management admin accounts
	protected.Get("/accounts", controller.ListAdminAccount) // Route to see all admin accounts
	protected.Get("/accounts/:id", controller.GetAccountByID) // Route to see admin account detail by acc_id
	protected.Put("/accounts/:id", controller.EditAdminAccount) // Route to edit password admin account by acc_id
	protected.Delete("/accounts/:id", controller.DeleteAdminAccount) // Route to delete admin account by acc_id

	// define routes for management certificate
	protected.Post("/certificate", nil)
}
