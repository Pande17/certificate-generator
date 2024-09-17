package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status": fiber.StatusBadRequest,
		"error": message,
		"timestamp": time.Now(),
		"data": nil,
	})
}

func Conflict(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		"status": fiber.StatusConflict,
		"error": message,
		"timestamp": time.Now(),
		"data": nil,
	})
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status": fiber.StatusInternalServerError,
		"error": message,
		"timestamp": time.Now(),
		"data": nil,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": fiber.StatusUnauthorized,
		"error": message,
		"timestamp": time.Now(),
		"data": nil,
	})
}

func Ok(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"message": message,
		"timestamp": time.Now(),
		"data": data,
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status": fiber.StatusNotFound,
		"error": message,
		"timestamp": time.Now(),
		"data": nil,
	})
}

func AlreadyDeleted(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status": fiber.StatusBadRequest,
		"error": message,
		"timestamp": time.Now(),
		"deleted_at": data,
	})
}

