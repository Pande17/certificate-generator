package routes

import (
	"certificate-generator/internal/handler/middleware"
	"certificate-generator/internal/handler/rest"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// function for setup routes
func RouteSetup(r *fiber.App) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "test",
		})
	})

	// CORS Middleware setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",                               // Replace with your frontend URL
		AllowMethods:     "GET,POST,PUT,DELETE",                                 // Allowed HTTP methods
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cookie", // Allowed headers
		ExposeHeaders:    "Authorization, Cookie",
		AllowCredentials: true,
	}))

	// Define a group routes for API
	api := r.Group("/api")

	// Define routes for authentication
	api.Post("/signup", middleware.AuditMiddleware("POST", "account"), rest.SignUp)                           // Route for signing up admin
	api.Post("/login", middleware.AuditMiddleware("POST", "account"), rest.Login)                             // Route for admin login
	api.Get("/validate", middleware.ValidateToken, rest.Validate)                                             // Route to check cookie from admin
	api.Post("/logout", middleware.ValidateToken, middleware.AuditMiddleware("POST", "account"), rest.Logout) // Route to logout from account

	// Define api routes
	// Every request to a path with the group "api" always checks the cookie
	// protected := api.Use(middleware.ValidateCookie)

	// Define routes for management admin accounts
	api.Get("/accounts", middleware.ValidateToken, middleware.AuditMiddleware("POST", "account"), rest.GetAdminAccount)             // Route to see all admin accounts
	api.Put("/accounts/:id", middleware.ValidateToken, middleware.AuditMiddleware("POST", "account"), rest.EditAdminAccount)        // Route to update password admin account by acc_id
	api.Delete("/accounts/:id", middleware.ValidateToken, middleware.AuditMiddleware("DELETE", "account"), rest.DeleteAdminAccount) // Route to delete admin account by acc_id

	// define routes for management competence
	api.Post("/competence", middleware.ValidateToken, middleware.AuditMiddleware("POST", "Competence"), rest.CreateKompetensi)         // route to create competence data
	api.Get("/competence", middleware.ValidateToken, middleware.AuditMiddleware("GET", "Competence"), rest.GetKompetensi)              // route to get all competence data
	api.Put("/competence/:id", middleware.ValidateToken, middleware.AuditMiddleware("PUT", "Competence"), rest.EditKompetensi)         // route to update competence data
	api.Delete("/competence/:id", middleware.ValidateToken, middleware.AuditMiddleware("DELETE", "Competence"), rest.DeleteKompetensi) // route to delete competence data

	// define routes for management certificate data
	api.Post("/certificate", middleware.ValidateToken, middleware.AuditMiddleware("POST", "Certificate"), rest.CreateCertificate)
	api.Get("/certificate", middleware.ValidateToken, middleware.AuditMiddleware("GET", "Certificate"), rest.GetAllCertificates)
	api.Get("/certificate/:id", middleware.ValidateToken, middleware.AuditMiddleware("GET", "Certificate"), rest.GetCertificateByID)
	api.Put("/certiticate/:id", middleware.ValidateToken, middleware.AuditMiddleware("PUT", "Certificate"), TEMPlate)
	api.Delete("/certificate/:id", middleware.ValidateToken, middleware.AuditMiddleware("DELETE", "Certificate"), rest.DeleteCertificate)

	r.Get("/assets/certificate/:id", rest.DownloadCertificate)
}

func TEMPlate(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON("nothing here yet")
}
