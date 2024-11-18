package middleware

import (
	"certificate-generator/database"
	"certificate-generator/model"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuditMiddleware(action, entity string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Call the next handler (perform the main action like GET, POST, etc.)
		err := c.Next()

		if err == nil {
			// Retrieve admin token from context
			adminClaims := c.Locals("admin")
			if adminClaims == nil {
				fmt.Println("No token found in context")
				return err
			}

			// Assert that adminToken is of type *jwt.Token
			// token, ok := adminToken.(*jwt.Token)
			// if !ok {
			// 	fmt.Println("Token is not of type *jwt.Token")
			// 	return err
			// }

			// Extract claims from token
			claims, ok := adminClaims.(jwt.MapClaims)
			if !ok {
				fmt.Println("Invalid token claims or token is not valid")
				return err
			}

			// Retrieve admin ID (subject) from claims
			adminIDHex, _ := claims["sub"].(string)
			adminID, _ := primitive.ObjectIDFromHex(adminIDHex)
			publicIP, _ := GetPublicIP()

			// Log audit action
			auditLog := model.AuditLog{
				ID:        primitive.NewObjectID(),
				AdminID:   adminID,
				Action:    action,
				Entity:    entity,
				EntityID:  c.Params("id", ""),
				Timestamp: time.Now(),
				IPAddress: publicIP,
			}

			collection := database.GetCollection("auditLog")
			_, err := collection.InsertOne(context.TODO(), auditLog)
			if err != nil {
				fmt.Println("Error saving audit log:", err)
			}
		}

		return err
	}
}
