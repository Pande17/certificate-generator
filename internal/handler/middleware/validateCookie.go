package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Middleware function to validate the JWT token from the cookie
func ValidateCookie(c *fiber.Ctx) error {
	// Retrieve the cookie named "Authorization"
	cookie := c.Cookies("Authorization")
	
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  "Unauthorized, please login",
			"status": fiber.StatusUnauthorized,
		})
	}

	// Parse the JWT token from the cookie
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key to validate the token
		return []byte(os.Getenv("SECRET")), nil
	})

	// If there is an error or the token is invalid
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  "Invalid token",
			"status": fiber.StatusUnauthorized,
		})
	}

	// Check if token is valid and not expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check for expiration time
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":  "Token expired",
					"status": fiber.StatusUnauthorized,
				})
			}
		}
		// Store the admin ID in the context for future use
		c.Locals("adminID", claims["sub"])
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  "Invalid token",
			"status": fiber.StatusUnauthorized,
		})
	}

	// Continue to the next handler
	return c.Next()
}
