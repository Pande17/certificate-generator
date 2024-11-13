package middleware

// import (
// 	"context"
// 	"fmt"
// 	"pkl/finalProject/certificate-generator/internal/database"
// 	"pkl/finalProject/certificate-generator/model"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v4"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// func AuditMiddleware(action, entity string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Call the next handler (perform the main action like GET, POST, PUT, DELETE)
// 		err := c.Next()

// 		if err == nil {
// 			// Retrieve admin token from context as a string
// 			adminTokenStr, ok := c.Locals("admin").(string)
// 			if !ok || adminTokenStr == "" {
// 				fmt.Println("No valid token found in context")
// 				return err
// 			}

// 			// Debug: Print the token to verify its format
// 			fmt.Printf("Retrieved adminTokenStr: %s\n", adminTokenStr)

// 			// Parse the token string to get *jwt.Token
// 			token, parseErr := jwt.Parse(adminTokenStr, func(token *jwt.Token) (interface{}, error) {
// 				// Replace `yourSecretKey` with the actual secret key used for signing JWTs
// 				return []byte("SECRET"), nil
// 			})

// 			if parseErr != nil {
// 				fmt.Println("Invalid token or parsing failed:", parseErr)
// 				return err
// 			}

// 			if !token.Valid {
// 				fmt.Println("Token is invalid")
// 				return err
// 			}

// 			// Extract claims from the token
// 			claims, ok := token.Claims.(jwt.MapClaims)
// 			if !ok {
// 				fmt.Println("Failed to extract claims from token")
// 				return err
// 			}

// 			// Retrieve admin ID (subject) from claims
// 			adminIDHex, _ := claims["sub"].(string)
// 			adminID, _ := primitive.ObjectIDFromHex(adminIDHex)
// 			publicIP, _ := GetPublicIP()

// 			collection := database.GetCollection("auditLog")

// 			// Log Audit action
// 			auditLog := model.AuditLog{
// 				ID:        primitive.NewObjectID(),
// 				AdminID:   adminID,
// 				Action:    action,
// 				Entity:    entity,
// 				EntityID:  c.Params("id", ""),
// 				Timestamp: time.Now(),
// 				IPAddress: publicIP,
// 			}

// 			_, err := collection.InsertOne(context.TODO(), auditLog)
// 			if err != nil {
// 				fmt.Println("Error saving Audit Log:", err)
// 			}
// 		}

// 		return err
// 	}
// }
