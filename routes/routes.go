package routes

import (
	"pkl/finalProject/certificate-generator/controller"
	"pkl/finalProject/certificate-generator/middleware"

	"github.com/gofiber/fiber/v2"
)

// function for setup routes
func RouteSetup(r *fiber.App) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "test",
		})
	})

	APIRoute(r)
	TemplateRoute(r)
}

// function for api routes
func APIRoute(r *fiber.App) {
	// Define a group routes for API
	api := r.Group("/api")

	// Define routes for authentication
	api.Post("/signup", controller.SignUp)                               // Route for signing up admin
	api.Post("/login", controller.Login)                                 // Route for admin login
	api.Get("/validate", middleware.ValidateCookie, controller.Validate) // Route to check cookie from admin

	// Define protected routes
	// Setiap request ke path dengan group "protected" selalu cek cookie
	protected := api.Use(middleware.ValidateCookie)

	// Define routes for authentication
	protected.Post("/logout", controller.Logout) // Route to logout from account

	// Define routes for management admin accounts
	protected.Get("/accounts", controller.GetAllAdminAccount)        // Route to see all admin accounts
	protected.Get("/accounts/:id", controller.GetAccountByID)        // Route to see admin account detail by acc_id
	protected.Put("/accounts/:id", controller.EditAdminAccount)      // Route to update password admin account by acc_id
	protected.Delete("/accounts/:id", controller.DeleteAdminAccount) // Route to delete admin account by acc_id

	// define routes for management competence
	protected.Post("/competence", controller.CreateKompetensi)       // route to create competence data
	protected.Put("/competence/:id", controller.EditKompetensi)      // route to update competence data
	protected.Delete("/competence/:id", controller.DeleteKompetensi) // route to delete competence data
	protected.Get("/competence", controller.GetAllKompetensi)        // route to get all competence data
	protected.Get("/competence/:id", controller.GetDetailKompetensi) // route to get competence data based on their id
	protected.Post("/qrcode", controller.CreateQRCode)               // route to create competence data
}

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
