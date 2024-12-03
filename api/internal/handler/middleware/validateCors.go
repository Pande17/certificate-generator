package middleware

import (
	"certificate-generator/internal/handler/rest"

	"github.com/gofiber/fiber/v2"
)

func CorsValidate(c *fiber.Ctx) error {
	if string(c.Request().Header.Peek("")) != "cors" {
		return rest.Unauthorized(c, "API aplikasi tidak dapat diakses secara langsung dari web.", "inaccessible through navigation. use web request instead.")
	}
	return nil
}
