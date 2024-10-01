package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func BadRequest(c *fiber.Ctx, message string) error {
	return status(c, fiber.StatusBadRequest, message, nil)
	// c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 	"status": fiber.StatusBadRequest,
	// 	"error": message,
	// 	"timestamp": time.Now(),
	// 	"data": nil,
	// })
}

func Conflict(c *fiber.Ctx, message string) error {
	return status(c, fiber.StatusConflict, message, nil)
	// c.Status(fiber.StatusConflict).JSON(fiber.Map{
	// 	"status": fiber.StatusConflict,
	// 	"error": message,
	// 	"timestamp": time.Now(),
	// 	"data": nil,
	// })
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return status(c, fiber.StatusInternalServerError, message, nil)
	// c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 	"status": fiber.StatusInternalServerError,
	// 	"error": message,
	// 	"timestamp": time.Now(),
	// 	"data": nil,
	// })
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return status(c, fiber.StatusUnauthorized, message, nil)
	// c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 	"status": fiber.StatusUnauthorized,
	// 	"error": message,
	// 	"timestamp": time.Now(),
	// 	"data": nil,
	// })
}

func Ok(c *fiber.Ctx, message string, data any) error {
	return status(c, fiber.StatusOK, message, data)
	// c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"status": fiber.StatusOK,
	// 	"message": message,
	// 	"timestamp": time.Now(),
	// 	"data": data,
	// })
}

func NotFound(c *fiber.Ctx, message string) error {
	return status(c, fiber.StatusNotFound, message, nil)
	// c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 	"status": fiber.StatusNotFound,
	// 	"error": message,
	// 	"timestamp": time.Now(),
	// 	"data": nil,
	// })
}

func AlreadyDeleted(c *fiber.Ctx, message string, data any) error {
	return status(c, 4000, message, data)
	// c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 	"status": fiber.StatusBadRequest,
	// 	"error": message,
	// 	"timestamp": time.Now(),
	// 	"deleted_at": data,
	// })
}

// for status Already Deleted, use 4000
func status(c *fiber.Ctx, stt int, message string, data any) error {
	fibMap := fiber.Map{
		"status":    stt,
		"timestamp": time.Now(),
		"data":      data,
		"message":   nil,
		"error":     message,
	}
	if stt == fiber.StatusOK {
		fibMap["message"] = message
		fibMap["error"] = nil
	}
	// idk which stt code to use for this
	if stt == 4000 {
		stt = 400
		fibMap["deleted_at"] = data
		fibMap["status"] = stt
		fibMap["data"] = nil
	}
	return c.Status(stt).JSON(fibMap)
}
