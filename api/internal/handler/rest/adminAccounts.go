package rest

import (
	"context"
	"fmt"
	"os"
	"time"

	"certificate-generator/api/database"
	model "certificate-generator/api/model"

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
		ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		AdminName     string             `json:"admin_name" bson:"admin_name"`
		AdminPassword string             `json:"admin_password" bson:"admin_password"`
	}

	// parse the request body
	if err := c.BodyParser(&adminReq); err != nil {
		return BadRequest(c, "Failed to read body", err.Error())
	}

	// validation to check if input username empty
	if adminReq.AdminName == "" {
		return BadRequest(c, "Admin Name cannot be empty", "Check admin name")
	}

	// validation to check if input password empty
	if adminReq.AdminPassword == "" {
		return BadRequest(c, "Password cannot be empty", "Check password empty")
	}

	// hashing the input password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(adminReq.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return BadRequest(c, "Failed to hashed password", "Hashing password signup")
	}

	// connect collection admin account in database
	adminCollection := database.GetCollection("adminAcc")

	// new variable to check the availability of the admin account name
	var existingAdmin model.AdminAccount

	// new variable to find admin account based on their "admin_name"
	filter := bson.M{"admin_name": adminReq.AdminName}

	// find admin account with same account name as input name
	err = adminCollection.FindOne(context.TODO(), filter).Decode(&existingAdmin)
	if err == nil {
		return Conflict(c, "Admin name already exists", "Check admin name")
	} else if err != mongo.ErrNoDocuments {
		return InternalServerError(c, "Error checking for existing admin name", "Check admin name")
	}

	// generate acc_id (incremental id)
	// nextAccID, err := generator.GetNextIncrementalID(adminCollection, "acc_id")
	// if err != nil {
	// 	return InternalServerError(c, "Failed to generate account ID", "Generate acc_id")
	// }

	// input data from req struct to struct "AdminAccount"
	admin := model.AdminAccount{
		ID:            primitive.NewObjectID(),
		AdminName:     adminReq.AdminName,
		AdminPassword: string(hash),
		Model: model.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// insert data from struct "AdminAccount" to collection in database MongoDB
	_, err = adminCollection.InsertOne(context.TODO(), admin)
	if err != nil {
		return InternalServerError(c, "Failed to create admin account", "Insert admin acc")
	}

	// return success
	return OK(c, "Admin account created successfully", admin)
}

// function to login to admin account
func Login(c *fiber.Ctx) error {
	// struct for the incoming request body
	var adminReq struct {
		AdminName     string `json:"admin_name" bson:"admin_name"`
		AdminPassword string `json:"admin_password" bson:"admin_password"`
	}

	// parse the request body
	if err := c.BodyParser(&adminReq); err != nil {
		return BadRequest(c, "Failed to read body", "Req body login")
	}

	// new variable to store admin login data
	var admin model.AdminAccount

	// connect collection admin account in database
	adminCollection := database.GetCollection("adminAcc")

	// find admin account with same account name as input name
	err := adminCollection.FindOne(context.TODO(), bson.M{"admin_name": adminReq.AdminName}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Unauthorized(c, "Invalid name or password", "Find admin name on login")
		}
		return InternalServerError(c, "Database error", "Database find admin name on login")
	}

	// Check if DeletedAt field already has a value (account has been deleted)
	if admin.DeletedAt != nil && !admin.DeletedAt.IsZero() {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This account has already been deleted", "Check deleted admin acc", admin.DeletedAt)
	}

	// hashing the input password with bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(adminReq.AdminPassword))
	if err != nil {
		return Unauthorized(c, "Invalid name or password", "Failed hashing password on Login")
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
		return InternalServerError(c, "Secret key not set", "Can not Retrieve SECRET key")
	}

	// use the secret key to sign the token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return BadRequest(c, "Failed to create Token", "Can not use SECRET key")
	}

	// set a cookie for admin
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		HTTPOnly: true,
	})

	// return success
	return OK(c, "Login successful", admin)
}

// Function to Validate checks if the user has a valid authentication cookie
func Validate(c *fiber.Ctx) error {
	// Retrieve the admin ID from the context set by the middleware
	adminID := c.Locals("admin")

	// Check if adminID is present (meaning the user is authenticated)
	if adminID == nil {
		return Unauthorized(c, "Unauthorized, please login", "Can not validate user")
	}

	// Return a success message along with the admin ID
	return OK(c, "User is authenticated", adminID)
}

// Function to logout
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})
	// return success
	return OK(c, "Logout Successful", nil)
}

// Function to edit data account
func EditAdminAccount(c *fiber.Ctx) error {
	// get acc_id from params
	idParam := c.Params("id")

	// converts acc_id to integer data type
	accID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Invalid ID format", "Can not convert params on Edit Admin")
	}

	// setup collection mongoDB
	adminCollection := database.GetCollection("adminAcc")

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"_id": accID}

	// variable to hold results
	var acc bson.M

	// Search for the account based on incremental ID
	if err := adminCollection.FindOne(c.Context(), filter).Decode(&acc); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Account not found", "Failed find account")
		}
		return InternalServerError(c, "Failed to fetch account", "Can not find account")
	}

	// Check if DeletedAt field already has a value
	if deletedAt, ok := acc["model"].(bson.M)["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This account has already been deleted", "Check deleted admin acc", deletedAt)
	}

	// Parsing req body to get new data
	var input struct {
		AdminName     string `json:"admin_name" bson:"admin_name"`
		AdminPassword string `json:"admin_password" bson:"admin_password"`
	}

	// handler if request body is invalid
	if err := c.BodyParser(&input); err != nil {
		return BadRequest(c, "Invalid request body", "Check req body")
	}
	// handler if admin password empty
	if input.AdminPassword == "" {
		return BadRequest(c, "Password cannot be empty", "Check empty password")
	}

	// hashing the input password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(input.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return BadRequest(c, "Failed to hash password", "")
	}

	// update fields in the database
	update := bson.M{
		"$set": bson.M{
			"admin_name":       input.AdminName,
			"admin_password":   string(hash),
			"model.updated_at": time.Now(),
		},
	}

	// update data in collection based on their "acc_id"
	_, err = adminCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return InternalServerError(c, "Failed to update account", "Can not Update admin acc")
	}

	// return success
	return OK(c, "Successfully updated account", update)
}

// Function for soft delete admin account
func DeleteAdminAccount(c *fiber.Ctx) error {
	// Get acc_id from params
	idParam := c.Params("id")

	// Converts acc_id to integer data type
	accID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Invalid ID format", "Can not convert params")
	}

	// connect to collection in mongoDB
	adminCollection := database.GetCollection("adminAcc")

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"_id": accID}

	// find admin account
	var adminAccount bson.M
	err = adminCollection.FindOne(context.TODO(), filter).Decode(&adminAccount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Account not found", "Can not find admin acc")
		}
		return InternalServerError(c, "Failed to fetch account details", "Can not find admin acc")
	}

	// Check if DeletedAt field already has a value
	if deletedAt, ok := adminAccount["model"].(bson.M)["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This account has already been deleted", "Check deleted admin acc", deletedAt)
	}

	// make update for input timestamp DeletedAt
	update := bson.M{"$set": bson.M{"model.deleted_at": time.Now()}}

	// update document in collection MongoDB
	result, err := adminCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return InternalServerError(c, "Failed to delete account", "Deleted admin acc")
	}

	// Check if the document is found and updated
	if result.MatchedCount == 0 {
		return NotFound(c, "Account not found", "Found admin acc")
	}

	// Respons success
	return OK(c, "Successfully deleted account", accID)
}

// func for get admin account
func GetAdminAccount(c *fiber.Ctx) error {
	if len(c.Queries()) == 0 {
		return getAllAdminAccount(c)
	}
	key := c.Query("type")
	val := c.Query("s")
	var value any
	if key == "id" {
		key = "_id"
		var err error
		if value, err = primitive.ObjectIDFromHex(val); err != nil {
			return BadRequest(c, "can't parse id", err.Error())
		}
	} else {
		value = val
	}
	return getOneAdminAccount(c, bson.M{key: value})
}

// Function to se all Admin Account
func getAllAdminAccount(c *fiber.Ctx) error {
	var results []bson.M

	adminCollection := database.GetCollection("adminAcc")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set the projection to return the required fields

	projection := bson.M{
		"_id":        1, // 1 to include the field, _id will be included by default
		"admin_name": 1, // 0 to exclude the field
		"created_at": 1,
		"updated_at": 1,
	}

	// find the projection
	cursor, err := adminCollection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "No Account Found", err.Error())
		}
		return InternalServerError(c, "Failed to fetch data", "Failed found admin acc")
	}
	defer cursor.Close(ctx)

	// decode each document and append it to results
	for cursor.Next(ctx) {
		var admin bson.M
		if err := cursor.Decode(&admin); err != nil {
			return InternalServerError(c, "Failed to decode data", "Can not decode data")
		}
		results = append(results, admin)
	}
	if err := cursor.Err(); err != nil {
		return InternalServerError(c, "Cursor error", "Cursor error")
	}

	return OK(c, "Success get all Admin Account data", results)
}

// function to get detail account by accID
func getOneAdminAccount(c *fiber.Ctx, filter bson.M) error {
	// connect to collection in MongoDB
	adminCollection := database.GetCollection("adminAcc")

	// Variable to hold search results
	var accountDetail bson.M

	// Get acc_id from params
	// idParam := c.Params("id")

	// parsing ObjectID to string
	// accID, err := primitive.ObjectIDFromHex(idParam)
	// if err != nil {
	// 	fmt.Printf("id on params: %v\n", idParam)
	// 	fmt.Printf("id received: %v\n", accID)
	// 	fmt.Printf("error: %v\n", err)
	// 	return BadRequest(c, "Invalid ID", "Cannot convert ID to ObjectID")
	// }

	// // make filter to find document based on id
	// filter := bson.M{"_id": accID}

	// Find a single document that matches the filter
	if err := adminCollection.FindOne(context.TODO(), filter).Decode(&accountDetail); err != nil {
		// If not found, return a 404 status.
		if err == mongo.ErrNoDocuments {
			fmt.Printf("error: %v\n", err)
			return NotFound(c, "Data not found 1", "Can not find account 2")
		}
		// If in server error, return status 500
		return InternalServerError(c, "Failed to retrieve data", "Server can't find account")
	}

	// check if document is already deleted
	if deletedAt, ok := accountDetail["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This account has already been deleted", "Check deleted admin acc", deletedAt)
	}

	// return success
	return OK(c, "Success get account data", accountDetail)
}
