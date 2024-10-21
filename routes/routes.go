package routes

import (
	"pkl/finalProject/certificate-generator/internal/handler/middleware"
	"pkl/finalProject/certificate-generator/internal/handler/rest"

	"github.com/gofiber/fiber/v2"
)

// function for setup routes
func RouteSetup(r *fiber.App) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "test",
		})
	})

	// Define a group routes for API
	api := r.Group("/api")

	// Define routes for authentication
	api.Post("/signup", rest.SignUp)                               // Route for signing up admin
	api.Post("/login", rest.Login)                                 // Route for admin login
	api.Get("/validate", middleware.ValidateCookie, rest.Validate) // Route to check cookie from admin
	api.Post("/testwk", rest.CreatePDF)

	// Define api routes
	// Every request to a path with the group "api" always checks the cookie
	// protected := api.Use(middleware.ValidateCookie)

	// Define routes for authentication
	api.Post("/logout", rest.Logout) // Route to logout from account

	// Define routes for management admin accounts
	api.Get("/accounts", rest.GetAllAdminAccount)        // Route to see all admin accounts
	api.Get("/accounts/:id", rest.GetAccountByID)        // Route to see admin account detail by acc_id
	api.Put("/accounts/:id", rest.EditAdminAccount)      // Route to update password admin account by acc_id
	api.Delete("/accounts/:id", rest.DeleteAdminAccount) // Route to delete admin account by acc_id

	// define routes for management competence
	api.Post("/competence", rest.CreateKompetensi)       // route to create competence data
	api.Get("/competence", rest.GetAllKompetensi)        // route to get all competence data
	api.Get("/competence/:id", rest.GetDetailKompetensi) // route to get competence data based on their id
	api.Put("/competence/:id", rest.EditKompetensi)      // route to update competence data
	api.Delete("/competence/:id", rest.DeleteKompetensi) // route to delete competence data

	// define routes for management certificate data
	api.Post("/certificate", rest.CreateCertificate)
	api.Get("/certificate", rest.GetAllCertificates)
	api.Get("/certificate/:id", rest.GetCertificateByID)
	api.Put("/certiticate/:id", rest.TEMPlate)
	api.Delete("/certificate/:id", rest.DeleteCertificate)

}
