package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"pkl/finalProject/certificate-generator/repository/config"
	dbmongo "pkl/finalProject/certificate-generator/repository/db_mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
    // Struct for the incoming request body
    var adminReq struct {
        ID            primitive.ObjectID `bson:"_id,omitempty"`
        AdminName     string             `json:"admin_name" bson:"admin_name"`
        AdminPassword string             `json:"admin_password" bson:"admin_password"`
    }

    // Parse the request body into the struct
    if err := c.BodyParser(&adminReq); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read body"})
    }

    // Hash the password using bcrypt
    hash, err := bcrypt.GenerateFromPassword([]byte(adminReq.AdminPassword), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to hash password"})
    }

    // Create the admin object to store in the database
    admin := dbmongo.AdminAccount{
        ID:            primitive.NewObjectID(), // Generate new ObjectID
        AdminName:     adminReq.AdminName,
        AdminPassword: string(hash),
        Model: dbmongo.Model{
            CreatedAt:     time.Now(),
            UpdatedAt:     time.Now(),
        },
    }

    // Insert the admin into MongoDB
    collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")
    _, err = collection.InsertOne(context.TODO(), admin)
    if err != nil {
        if mongo.IsDuplicateKeyError(err) {
            return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Admin name already exists"})
        }
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create admin"})
    }

    // Respond with the created admin information
    return c.Status(http.StatusOK).
        JSON(fiber.Map{
            "message": "Admin created successfully",
            "admin":   admin,
        })
}


func Login(c *fiber.Ctx) error {
	// Get the name and password from the request body
	var adminReq struct {
		AdminName     string `json:"admin_name" bson:"admin_name"`
		AdminPassword string `json:"admin_password" bson:"admin_password"`
	}

	if err := c.BodyParser(&adminReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read body"})
	}

	// Look up the requested user in MongoDB by both admin_name and admin_password
	var admin dbmongo.AdminAccount
	collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	// Search by admin_name only, we will verify the password separately
	err := collection.FindOne(context.TODO(), bson.M{"admin_name": adminReq.AdminName}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid name or password"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Compare the sent password with the saved user password hash
	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(adminReq.AdminPassword))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid name or password"})
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as string using a secret key
	secret := os.Getenv("SECRET")
	if secret == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Secret key not set"})
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create Token"})
	}

	// Send the token back as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		HTTPOnly: true,
	})

	return c.Status(http.StatusOK).
		JSON(fiber.Map{
			"admin":   admin,
			"message": "Login successful",
		})
}


func Validate(c *fiber.Ctx) error {
	admin, ok := c.Locals("Admin").(*dbmongo.AdminAccount)
	if !ok {
		fmt.Println("admin not found!")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "E Unauthorized"})
	}

	fmt.Println("Admin found in context:", admin)
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": admin})
}

