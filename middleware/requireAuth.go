package middleware

import (
	"context"
	"fmt"
	"os"
	"pkl/finalProject/certificate-generator/repository/config"
	dbmongo "pkl/finalProject/certificate-generator/repository/db_mongo"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RequireAuth(c *fiber.Ctx) error {
	// Get the cookie from request
	tokenString := c.Cookies("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "A Unauthorized"})
	}

	// Decode/validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "B Unauthorized"})
	}

	// Check the claims in the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "C Unauthorized"})
	}

	// Check if the token is expired
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
	}

	// Convert the 'sub' claim to ObjectID
	adminID, err := primitive.ObjectIDFromHex(claims["sub"].(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token payload"})
	}

	// Find the admin by token subject
	var admin dbmongo.AdminAccount
	collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")
	err = collection.FindOne(context.TODO(), bson.M{"_id": adminID}).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "D Unauthorized"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Attach admin to request context
	c.Locals("admin", admin)

	// Continue to the next middleware or handler
	return c.Next()
}
