package controller

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"pkl/finalProject/certificate-generator/repository/config"
	dbmongo "pkl/finalProject/certificate-generator/repository/db_mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Function to make new admin account
func SignUp(c *fiber.Ctx) error {
	// Struct for the incoming request body	
    var adminReq struct {
        ID            primitive.ObjectID `bson:"_id,omitempty"`
        AccID         uint64             `bson:"acc_id"`
        AdminName     string             `json:"admin_name" bson:"admin_name"`
        AdminPassword string             `json:"admin_password" bson:"admin_password"`
    }

	// parse the request body
    if err := c.BodyParser(&adminReq); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to read body",
            "status": fiber.StatusBadRequest,
        })
    }

    // Trim whitespace from input
    adminReq.AdminName = strings.TrimSpace(adminReq.AdminName)
    adminReq.AdminPassword = strings.TrimSpace(adminReq.AdminPassword)

    // validation to check if input password empty
    if adminReq.AdminPassword == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Password cannot be empty",
            "status": fiber.StatusBadRequest,
        })
    }

	// hashing the input password with bcrypt
    hash, err := bcrypt.GenerateFromPassword([]byte(adminReq.AdminPassword), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to hash password",
            "status": fiber.StatusBadRequest,
        })
    }

	// connect collection admin account in database
    collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	// new variable to check the availability of the admin account name
    var existingAdmin dbmongo.AdminAccount

	// new variable to find admin account based on their "admin_name"
    filter := bson.M{"admin_name": adminReq.AdminName}

	// find admin account with same account name as input name
    err = collection.FindOne(context.TODO(), filter).Decode(&existingAdmin)
    if err == nil {
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "error": "Admin name already exists",
            "status": fiber.StatusConflict,
        })
    } else if err != mongo.ErrNoDocuments {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error checking for existing admin name",
            "status": fiber.StatusInternalServerError,
        })
    }

	// generate acc_id (incremental id)
    nextAccID, err := GetNextAccID(collection)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to generate account ID",
            "status": fiber.StatusInternalServerError,
        })
    }

	// input data from req struct to struct "AdminAccount"
    admin := dbmongo.AdminAccount{
        ID:            primitive.NewObjectID(),
        AccID:         nextAccID,
        AdminName:     adminReq.AdminName,
        AdminPassword: string(hash),
        Model: dbmongo.Model{
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            DeletedAt: nil,
        },
    }

	// insert data from struct "AdminAccount" to collection in database MongoDB
    _, err = collection.InsertOne(context.TODO(), admin)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create admin",
            "status": fiber.StatusInternalServerError,
        })
    }

	// return success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Admin account created successfully",
        "status":  fiber.StatusOK,
        "admin":   admin,
    })
}

// function to login to admin account
func Login(c *fiber.Ctx) error {
	// struct for the incoming request body
    var adminReq struct {
        AdminName     string `json:"admin_name"`
        AdminPassword string `json:"admin_password"`
    }

	// parse the request body
    if err := c.BodyParser(&adminReq); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to read body",
            "status": fiber.StatusBadRequest,
        })
    }

    // Trim whitespace dari input
    adminReq.AdminName = strings.TrimSpace(adminReq.AdminName)
    adminReq.AdminPassword = strings.TrimSpace(adminReq.AdminPassword)

	// new variable to store admin login data
    var admin dbmongo.AdminAccount

	// connect collection admin account in database
    collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	// find admin account with same account name as input name
    err := collection.FindOne(context.TODO(), bson.M{"admin_name": adminReq.AdminName}).Decode(&admin)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid name or password",
                "status": fiber.StatusUnauthorized,
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Database error",
            "status": fiber.StatusInternalServerError,
        })
    }

	// hashing the input password with bcrypt
    err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(adminReq.AdminPassword))
    if err != nil {
       	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid name or password",
            "status": fiber.StatusUnauthorized,
        })
    }

	// generate token JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": admin.ID.Hex(),
        "exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
    })

	// retrieve the "SECRET" environment variable, which is used as the secret key for signing the JWT
    secret := os.Getenv("SECRET")
	// check if the secret key is not set (empty string)
    if secret == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Secret key not set",
            "status": fiber.StatusInternalServerError,
        })
    }

	// use the secret key to sign the token
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to create Token",
            "status": fiber.StatusBadRequest,
        })
    }

	// set a cookie for admin
    c.Cookie(&fiber.Cookie{
        Name:     "Authorization",
        Value:    tokenString,
        Expires:  time.Now().Add(24 * time.Hour * 30),
        HTTPOnly: true,
    })

	// return success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Login successful",
        "status":  fiber.StatusOK,
        "admin":   admin,
    })
}

// Function to Validate checks if the user has a valid authentication cookie
func Validate(c *fiber.Ctx) error {
	// Retrieve the admin ID from the context set by the middleware
	adminID := c.Locals("adminID")

	// Check if adminID is present (meaning the user is authenticated)
	if adminID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
			"error": "Unauthorized, please login",
		})
	}

	// Return a success message along with the admin ID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"adminID": adminID,
		// "account": ,
		"message": "User is authenticated",
		"status": fiber.StatusOK,
	})
}

// Function to logout
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name: "Authorization",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout Successful",
		"status": fiber.StatusOK,
	})
}

// Function to se all Admin Account
func ListAdminAccount(c *fiber.Ctx) error {
	var results []dbmongo.AdminAccount

	collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No Documents Found!",
				"status": fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data",
			"status": fiber.StatusInternalServerError,	
		})
	}	
	defer cursor.Close(ctx)


	for cursor.Next(ctx) {
		var admin dbmongo.AdminAccount
		if err := cursor.Decode(&admin); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode data",
				"status": fiber.StatusNotFound,	
			})
		}
		results = append(results, admin)
	}
	if err := cursor.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cursor error",
			"status": fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success get all data",
		"status": fiber.StatusOK,
		"data": results,
	})
}

// function to get detail account by accID
func GetAccountByID(c *fiber.Ctx) error {
	// Get acc_id from params
	idParam := c.Params("id")

	// parsing acc_id to integer type data
	accID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
			"status": fiber.StatusBadRequest,
		})
	}

	// connect to collection in mongoDB
	collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"acc_id": accID}

	// Variable to hold search results
	var accountDetail bson.M
	
	// Find a single document that matches the filter
	err = collection.FindOne(context.TODO(), filter).Decode(&accountDetail)
	if err != nil {
		// If not found, return a 404 status.
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  fiber.StatusNotFound,
				"message": "Data not found",
				"data":    nil,
		})
	}

	// If in server error, return status 500
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  fiber.StatusInternalServerError,
		"error": "Failed to retrieve data",
	})
}
	// return success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": fiber.StatusOK,
			"Message": "Data found",
			"data": accountDetail,
		})
}

// Function to edit data account
func EditAdminAccount(c *fiber.Ctx) error {
	// get acc_id from params
	idParam := c.Params("id")

	// converts acc_id to integer data type
	accID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
			"status": fiber.StatusBadRequest,
		})
	}

	// setup collection mongoDB
    collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"acc_id": accID}

	// variable to hold results
	var acc bson.M

    // Search for the account based on incremental ID
    if err := collection.FindOne(c.Context(), filter).Decode(&acc); err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error":  "Account ID not found",
            "status": fiber.StatusNotFound,
        })
    }

    // Parsing req body to get new data
    var input struct {
        AdminName string `json:"admin_name" bson:"admin_name"`
        AdminPassword  string `json:"admin_password" bson:"admin_password"`
    }

	// handler if request body is invalid
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":  "Invalid request body",
            "status": fiber.StatusBadRequest,
        })
    }

	// handler if admin password empty
	if input.AdminPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password cannot be empty",
			"status": fiber.StatusBadRequest,
		})
	}

	// hashing the input password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(input.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to hash password",
			"status": fiber.StatusBadRequest,
		})
	}

    // update fields in the database
    update := bson.M{
		"$set": bson.M{
			"admin_name": input.AdminName,
			"admin_password": string(hash),
			"model.updated_at": time.Now(),
		},
	}

	// update data in collection based on their "acc_id"
    _, err = collection.UpdateOne(c.Context(), filter, update)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":  "Failed to update account",
            "status": fiber.StatusInternalServerError,
        })
    }

	// return success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Successfully updated account",
        "status":  fiber.StatusOK,
        "new_data": update,
    })
}

// Function for soft delete admin account
func DeleteAdminAccount(c *fiber.Ctx) error {
	// Get acc_id from params
	idParam := c.Params("id")

	// Converts acc_id to integer data type
	accID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid ID format",
			"status": fiber.StatusBadRequest,
		})
	}

	// connect to collection in mongoDB
	collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"acc_id": accID}

	// make update for input timestamp DeletedAt
	update := bson.M{"$set": bson.M{"model.deleted_at": time.Now()}}

	// update document in collection MongoDB
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to soft delete account",
			"status": fiber.StatusInternalServerError,
		})
	}

	// Check if the document is found and updated
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  "Account not found",
			"status": fiber.StatusNotFound,
		})
	}

	// Respons success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted account",
		"status":  fiber.StatusOK,
	})
}

// Function for generate ID incremental for admin account
func GetNextAccID(adminCollection *mongo.Collection) (int64, error) {
    // Define a filter to find the maximum AccID
    opts := options.FindOne().SetSort(bson.D{{"acc_id", -1}}) // Sort by acc_id descending

    var lastAdmin dbmongo.AdminAccount
	var ctx = context.Background() // Define the context

    // Retrieve the last inserted admin account
    err := adminCollection.FindOne(ctx, bson.M{}, opts).Decode(&lastAdmin)
    if err != nil && err != mongo.ErrNoDocuments {
        return 0, fmt.Errorf("failed to find the last admin account: %v", err)
    }

    // If no documents exist, start from 1
    if err == mongo.ErrNoDocuments {
        return 1, nil
    }

    // Increment the last AccID by 1
    return int64(lastAdmin.AccID)+1, nil
}