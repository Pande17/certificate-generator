package rest

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// helper function to build response
func jsonResponse(c *fiber.Ctx, statusCode int, message string, errLocate string, data any, deletedAt any) error {
	response := fiber.Map{
		"status":    statusCode,
		"message":   message,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// only include error_location if it's not empty
	if errLocate != "" {
		response["error_location"] = errLocate
	}

	// only include data if it's not nil
	if data != nil {
		response["data"] = data
	}

	// only include delete_at if its's not empty
	if deletedAt != nil {
		response["deleted_at"] = deletedAt
	}

	addAllowOrigin(c)

	return c.Status(statusCode).JSON(response)
}

// code 200
func OK(c *fiber.Ctx, message string, data any) error {
	return jsonResponse(c, fiber.StatusOK, message, "", data, nil)
}

// code 400
func BadRequest(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusBadRequest, message, errLocate, nil, nil)
}

// code 401
func Unauthorized(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusUnauthorized, message, errLocate, nil, nil)
}

// code 404
func NotFound(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusNotFound, message, errLocate, nil, nil)
}

// code 404 for deleted document
func AlreadyDeleted(c *fiber.Ctx, message string, errLocate string, deletedAt any) error {
	return jsonResponse(c, fiber.StatusNotFound, message, errLocate, nil, deletedAt)
}

// code 409
func Conflict(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusConflict, message, errLocate, nil, nil)
}

func addAllowOrigin(c *fiber.Ctx) {
	for _, url := range strings.Split(os.Getenv("CERTIF_GEN_FRONTEND"), ",") {
		if url == string(c.Request().Header.Peek("Origin")) {
			c.Response().Header.Set("Access-Control-Allow-Origin", url)
			break
		}
	}
}
