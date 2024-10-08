package routes

import (
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
