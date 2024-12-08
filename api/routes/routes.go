package routes

import (
	"certificate-generator/internal/handler/middleware"
	"certificate-generator/internal/handler/rest"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// function for setup routes
func RouteSetup(r *fiber.App) {
	// CORS Middleware setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CERTIF_GEN_FRONTEND"),
		AllowMethods:     "GET,POST,PUT,DELETE",                                                                      // Allowed HTTP methods
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cookie, authToken, Bearer, X-Requested-With", // Allowed headers
		ExposeHeaders:    "Content-Type, Authorization, Cookie, authToken, Bearer. Accept",
		AllowCredentials: true,
	}))

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "test",
		})
	})

	// Define a group routes for API
	api := r.Group("/api")

	api.Use(middleware.CorsValidate)

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "test",
		})
	})

	// Define routes for authentication
	api.Post("/signup", middleware.AuditMiddleware("SignUp"), rest.SignUp)
	api.Post("/login", middleware.AuditMiddleware("Login"), rest.Login)
	api.Get("/validate", middleware.ValidateToken, rest.Validate)
	api.Post("/logout", middleware.AuditMiddleware("LogOut"), rest.Logout)

	// Define routes for managing admin accounts
	api.Get("/accounts", middleware.ValidateToken, middleware.AuditMiddleware("Account"), rest.GetAdminAccount)
	api.Put("/accounts/:id", middleware.ValidateToken, middleware.AuditMiddleware("Account"), rest.EditAdminAccount)
	api.Delete("/accounts/:id", middleware.ValidateToken, middleware.AuditMiddleware("Account"), rest.DeleteAdminAccount)

	// Define routes for managing competence data
	api.Post("/competence", middleware.ValidateToken, middleware.AuditMiddleware("Competence"), rest.CreateKompetensi)
	api.Get("/competence", middleware.ValidateToken, middleware.AuditMiddleware("Competence"), rest.GetAllKompetensi)
	api.Get("/competence/:id", middleware.ValidateToken, middleware.AuditMiddleware("Competence"), rest.GetKompetensiByID)
	api.Get("/competence/:type/:id", middleware.ValidateToken, middleware.AuditMiddleware("Competence"), rest.GetKompetensiByID)
	api.Put("/competence/:id", middleware.ValidateToken, middleware.AuditMiddleware("Competence"), rest.EditKompetensi)
	api.Delete("/competence/:id", middleware.ValidateToken, middleware.AuditMiddleware("Competence"), rest.DeleteKompetensi)

	// Define routes for managing certificate data
	api.Post("/certificate", middleware.ValidateToken, middleware.AuditMiddleware("Certificate"), rest.CreateCertificate)
	api.Get("/certificate", middleware.ValidateToken, middleware.AuditMiddleware("Certificate"), rest.GetAllCertificates)
	api.Get("/certificate/:id", middleware.ValidateToken, middleware.AuditMiddleware("Certificate"), rest.GetCertificateByID)
	api.Get("/certificate/:type/:id", middleware.ValidateToken, middleware.AuditMiddleware("Certificate"), rest.GetCertificateByID)
	api.Put("/certificate/:id", middleware.ValidateToken, middleware.AuditMiddleware("Certificate"), rest.EditCertificate)
	api.Delete("/certificate/:id", middleware.ValidateToken, middleware.AuditMiddleware("Certificate"), rest.DeleteCertificate)
	api.Get("/certificate/download/:id/:type", rest.DownloadCertificate, rest.GetCertificateByID)

	// Define routes for managing signature configuration data
	api.Post("/signature", middleware.ValidateToken, middleware.AuditMiddleware("Signature"), rest.CreateSignature)
	api.Get("/signature", middleware.ValidateToken, middleware.AuditMiddleware("Signature"), rest.GetAllSignature)
	api.Get("/signature/:id", middleware.ValidateToken, middleware.AuditMiddleware("Signature"), rest.GetSignatureByID)
	api.Get("/signature/:type/:id", middleware.ValidateToken, middleware.AuditMiddleware("Signature"), rest.GetSignatureByID)
	api.Put("/signature/:id", middleware.ValidateToken, middleware.AuditMiddleware("Signature"), rest.EditSignature)
	api.Delete("/signature/:id", middleware.ValidateToken, middleware.AuditMiddleware("Signature"), rest.DeleteSignature)
}

func TEMPlate(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON("nothing here yet")
}
