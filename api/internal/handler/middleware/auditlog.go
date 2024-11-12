package middleware

import (
	"context"
	"fmt"
	"pkl/finalProject/certificate-generator/internal/database"
	"pkl/finalProject/certificate-generator/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuditMiddleware(action, entity string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// call the next handler (perform the main action like GET, POST, PUT, DELETE)
		err := c.Next()

		if err == nil {
			// retrieve admin token from context
			adminToken := c.Locals("adminID")
			fmt.Printf("Type of adminToken: %T\n", adminToken) // Ini akan membantu verifikasi tipe
			if adminToken == nil {
				fmt.Println("No token found in context")
				fmt.Printf("error: %v\n", err)
				return err
			}

			// assert that adminToken is *jwt.token
			token, ok := adminToken.(*jwt.Token)
			if !ok {
				fmt.Println("Invalid token claims or token is not valid coy")
				return err
			}

			// extract claims for token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				fmt.Println("Invalid token claims or token is not valid cuyy")
				return err
			}

			// retrieve admin ID (subject) from claims
			adminIDHex, _ := claims["sub"].(string)
			adminID, _ := primitive.ObjectIDFromHex(adminIDHex)
			publicIP, _ := GetPublicIP()

			collection := database.GetCollection("auditLog")

			// Log Audit action
			auditLog := model.AuditLog{
				ID:        primitive.NewObjectID(),
				AdminID:   adminID,
				Action:    action,
				Entity:    entity,
				EntityID:  c.Params("id", ""),
				Timestamp: time.Now(),
				IPAddress: publicIP,
			}

			_, err := collection.InsertOne(context.TODO(), auditLog)
			if err != nil {
				fmt.Println("Error saving Audit Log: ", err)
			}
		}

		return err
	}
}
