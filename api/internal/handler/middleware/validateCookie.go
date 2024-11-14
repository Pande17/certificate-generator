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

// Middleware function to validate the JWT token from the cookie
// func ValidateCookie(c *fiber.Ctx) error {
// 	// Retrieve the cookie named "Authorization"
// 	cookie := c.Cookies("Authorization")

// 	if cookie == "" {
// 		fmt.Println("Cookie token:", cookie)
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error":  "Unauthorized, please login",
// 			"status": fiber.StatusUnauthorized,
// 		})
// 	}

// 	// Parse the JWT token from the cookie
// 	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
// 		// Ensure the signing method is HMAC
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		// Return the secret key to validate the token
// 		return []byte(os.Getenv("SECRET")), nil
// 	})

// 	// If there is an error or the token is invalid
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error":  "Invalid token",
// 			"status": fiber.StatusUnauthorized,
// 		})
// 	}

// 	// Check if token is valid and not expired
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		// Check for expiration time
// 		if exp, ok := claims["exp"].(float64); ok {
// 			if time.Unix(int64(exp), 0).Before(time.Now()) {
// 				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 					"error":  "Token expired",
// 					"status": fiber.StatusUnauthorized,
// 				})
// 			}
// 		}
// 		// Store the admin ID in the context for future use
// 		c.Locals("admin", claims["sub"])
// 	} else {
// 		fmt.Println("Invalid token or claims:", err)
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error":  "Invalid token 2",
// 			"status": fiber.StatusUnauthorized,
// 		})
// 	}

// 	// Continue to the next handler
// 	return c.Next()
// }

// ValidateToken middleware
func ValidateToken(c *fiber.Ctx) error {
	// Retrieve the token from Authorization header or cookies
	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		tokenString = strings.TrimPrefix(c.Cookies("Authorization"), "Bearer ")
	}

	fmt.Println("Token:", tokenString)

	// If token is missing, return unauthorized error
	if tokenString == "" {
		return rest.Unauthorized(c, "Unauthorized, please login", "Unauthorized, please login")
	}

	// Retrieve secret key from environment variables
	secretKey := os.Getenv("SECRET")
	if secretKey == "" {
		return rest.Unauthorized(c, "Server error", "Server error")
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
		fmt.Println("Error parsing token:", err)
		return rest.Unauthorized(c, "Invalid Token", "Invalid Token")
	}

	// Log the type of token returned
	fmt.Printf("Parsed token: %T\n", token) // Log the type of the token object

	// Ensure the token is valid
	if !token.Valid {
		fmt.Println("Invalid token")
		return rest.Unauthorized(c, "Invalid Token", "Invalid Token")
	}

	// Check for expiration (if the exp claim exists)
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return rest.Unauthorized(c, "Expired Token", "Expired Token")
		}
	}

	// Log the claims for debugging purposes
	fmt.Println("Token claims:", claims)

	// Store the entire claims map in context for later use
	c.Locals("admin", claims) // Store all claims in context

	// Proceed to the next middleware or handler
	return c.Next()
}
