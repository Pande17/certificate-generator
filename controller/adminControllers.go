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
		return BadRequest(c, "Failed to read body")
    }

    // Trim whitespace from input
    adminReq.AdminName = strings.TrimSpace(adminReq.AdminName)
    adminReq.AdminPassword = strings.TrimSpace(adminReq.AdminPassword)

    // validation to check if input password empty
    if adminReq.AdminPassword == "" {
        return BadRequest(c, "Password cannot be empty")
    }

	// hashing the input password with bcrypt
    hash, err := bcrypt.GenerateFromPassword([]byte(adminReq.AdminPassword), bcrypt.DefaultCost)
    if err != nil {
        return BadRequest(c, "Failed to hashed password")
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
		return Conflict(c, "Admin name already exists")
    } else if err != mongo.ErrNoDocuments {
		return InternalServerError(c, "Error checking for existing admin name")
    }

	// generate acc_id (incremental id)
    nextAccID, err := GetNextAccID(collection)
    if err != nil {
		return InternalServerError(c, "Failed to generate account ID")
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
        return InternalServerError(c, "Failed to create admin account")
    }

	// return success
    return Ok(c, "Admin account created successfully", admin)
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
        return BadRequest(c, "Failed to read body")
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
            return Unauthorized(c, "Invalid name or password")
        }
        return InternalServerError(c, "Database error")
    }

	// Check if DeletedAt field already has a value (account has been deleted)
	if admin.DeletedAt != nil && !admin.DeletedAt.IsZero() {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This account has already been deleted", admin.DeletedAt)
	}

	// hashing the input password with bcrypt
    err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(adminReq.AdminPassword))
    if err != nil {
       	return Unauthorized(c, "Invalid name or password")
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
        return InternalServerError(c, "Secret key not set")
    }

	// use the secret key to sign the token
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        return BadRequest(c, "Failed to create Token")
    }

	// set a cookie for admin
    c.Cookie(&fiber.Cookie{
        Name:     "Authorization",
        Value:    tokenString,
        Expires:  time.Now().Add(24 * time.Hour * 30),
        HTTPOnly: true,
    })

	// return success
    return Ok(c, "Login successful", admin)
}

// Function to Validate checks if the user has a valid authentication cookie
func Validate(c *fiber.Ctx) error {
	// Retrieve the admin ID from the context set by the middleware
	adminID := c.Locals("adminID")

	// Check if adminID is present (meaning the user is authenticated)
	if adminID == nil {
		return Unauthorized(c, "Unauthorized, please login")
	}

	// Return a success message along with the admin ID
	return Ok(c, "User is authenticated", adminID)
}

// Function to logout
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name: "Authorization",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})
	// return success
	return Ok(c, "Logout Successful", nil)
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
			return NotFound(c, "No Documents Found")
		}
		return InternalServerError(c, "Failed to fetch data")
	}	
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var admin dbmongo.AdminAccount
		if err := cursor.Decode(&admin); err != nil {
			return InternalServerError(c, "Failed to decode data")
		}
		results = append(results, admin)
	}
	if err := cursor.Err(); err != nil {
		return InternalServerError(c, "Cursor error")
	}

	return Ok(c, "Success get all data", results)
}
// function to get detail account by accID
func GetAccountByID(c *fiber.Ctx) error {
	// Get acc_id from params
	idParam := c.Params("id")

	// parsing acc_id to integer type data
	accID, err := strconv.Atoi(idParam)
	if err != nil {
		return BadRequest(c, "Invalid ID format")
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
			return NotFound(c, "Data not found")
	}
	// If in server error, return status 500
	return InternalServerError(c, "Failed to retrieve data")
	}

	// Check if DeletedAt field already has a value
	if deletedAt, ok := accountDetail["model"].(bson.M)["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This account has already been deleted", deletedAt)
	}

	// return success
	return Ok(c, "Success get data", accountDetail)
}

// Function to edit data account
func EditAdminAccount(c *fiber.Ctx) error {
	// get acc_id from params
	idParam := c.Params("id")

	// converts acc_id to integer data type
	accID, err := strconv.Atoi(idParam)
	if err != nil {
		return BadRequest(c, "Invalid ID format")
	}

	// setup collection mongoDB
    collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"acc_id": accID}

	// variable to hold results
	var acc bson.M

    // Search for the account based on incremental ID
    if err := collection.FindOne(c.Context(), filter).Decode(&acc); err != nil {
        if err == mongo.ErrNoDocuments {
			return NotFound(c, "Account not found")
		}
		return InternalServerError(c, "Failed to fetch account")
    }

	// Check if DeletedAt field already has a value
	if deletedAt, ok := acc["model"].(bson.M)["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This account has already been deleted", deletedAt)
	}

    // Parsing req body to get new data
    var input struct {
        AdminName string `json:"admin_name" bson:"admin_name"`
        AdminPassword  string `json:"admin_password" bson:"admin_password"`
    }

	// handler if request body is invalid
    if err := c.BodyParser(&input); err != nil {
        return BadRequest(c, "Invalid request body")
	}
	// handler if admin password empty
	if input.AdminPassword == "" {
		return BadRequest(c, "Password cannot be empty")
	}

	// hashing the input password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(input.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return BadRequest(c, "Failed to hash password")
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
        return InternalServerError(c, "Failed to update account")
    }

	// return success
    return Ok(c, "Successfully updated account", update)
}

// Function for soft delete admin account
func DeleteAdminAccount(c *fiber.Ctx) error {
	// Get acc_id from params
	idParam := c.Params("id")

	// Converts acc_id to integer data type
	accID, err := strconv.Atoi(idParam)
	if err != nil {
		return BadRequest(c, "Invalid ID format")
	}

	// connect to collection in mongoDB
	collection := config.MongoClient.Database("certificate-generator").Collection("adminAcc")

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"acc_id": accID}


	// Check if account is already deleted
	var adminAccount bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&adminAccount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Account not found")
		}
		return InternalServerError(c, "Failed to fetch account details")
	}

	// Check if DeletedAt field already has a value
	if deletedAt, ok := adminAccount["model"].(bson.M)["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This account has already been deleted", deletedAt)
	}

	// make update for input timestamp DeletedAt
	update := bson.M{"$set": bson.M{"model.deleted_at": time.Now()}}

	// update document in collection MongoDB
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return InternalServerError(c, "Failed to delete account")
	}

	// Check if the document is found and updated
	if result.MatchedCount == 0 {
		return NotFound(c, "Account not found")
	}

	// Respons success
	return Ok(c, "Successfully deleted account", accID)
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