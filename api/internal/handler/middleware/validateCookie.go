package middleware

import (
	"certificate-generator/internal/handler/rest"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// ValidateToken middleware
func ValidateToken(c *fiber.Ctx) error {
	// Retrieve the token from authToken header or cookies
	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		tokenString = c.Cookies("authToken") // Use a specific cookie name
	}

	// If token is missing, return unauthorized error
	if tokenString == "" {
		return rest.Unauthorized(c, "Mohon login terlebih dahulu", "Token tidak ditemukan")
	}

	// Retrieve secret key from environment variables
	secretKey := os.Getenv("SECRET")
	if secretKey == "" {
		return rest.Conflict(c, "Server error", "Kunci rahasia tidak ditemukan")
	}

	// Initialize claims as MapClaims to store all claims
	claims := jwt.MapClaims{}

	// Parse the JWT token with the claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	// Check for errors during parsing
	if err != nil {
		return rest.Unauthorized(c, "Invalid Token", "Token tidak valid")
	}

	// Ensure the token is valid
	if !token.Valid {
		return rest.Unauthorized(c, "Invalid Token", "Token tidak valid")
	}

	// Check for expiration (if the exp claim exists)
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return rest.Unauthorized(c, "Expired Token", "Token telah kedaluwarsa")
		}
	}

	// Store the entire claims map in context for later use
	c.Locals("admin", claims) // Store all claims in context

	fmt.Println("Token from Authorization header:", tokenString)
	fmt.Println("Token from cookies:", c.Cookies("authToken"))

	// Proceed to the next middleware or handler
	return c.Next()
}
