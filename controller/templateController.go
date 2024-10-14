package controller

import "github.com/gofiber/fiber/v2"

func TEMPlate(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON("nothing here yet")
}
