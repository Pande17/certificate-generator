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

// ValidateToken middleware to check the validity of the JWT token
func ValidateToken(c *fiber.Ctx) error {
	// Retrieve the token from the Authorization header
	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		// If the token is not found in the header, check cookies
		tokenString = c.Cookies("authToken") // Use a specific cookie name
	}

	// If token is missing, return unauthorized error
	if tokenString == "" {
		return rest.Unauthorized(c, "Silakan login untuk mengakses fitur ini.", "Token tidak ditemukan.")
	}

	// Retrieve secret key from environment variables
	secretKey := os.Getenv("SECRET")
	if secretKey == "" {
		return rest.Conflict(c, "Terjadi kesalahan pada server. Kunci rahasia tidak ditemukan.", "Silakan hubungi admin.")
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
		return rest.Unauthorized(c, "Token yang Anda masukkan tidak valid. Silakan coba lagi.", "Token tidak valid.")
	}

	// Ensure the token is valid
	if !token.Valid {
		return rest.Unauthorized(c, "Token tidak valid. Pastikan Anda menggunakan token yang benar.", "Token tidak valid.")
	}

	// Check for expiration (if the exp claim exists)
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return rest.Unauthorized(c, "Token Anda telah kedaluwarsa. Silakan login kembali.", "Token telah kedaluwarsa.")
		}
	}

	// Store the entire claims map in context for later use
	c.Locals("admin", claims) // Store all claims in context

	// Proceed to the next middleware or handler
	return c.Next()
}
