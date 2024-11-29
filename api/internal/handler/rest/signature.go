package rest

import (
	"certificate-generator/database"
	"certificate-generator/model"
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get the "signature" collection from the database
var collectionSignature = database.GetCollection("signature")

// Function to create a new signature
func CreateSignature(c *fiber.Ctx) error {
	// Struct for the incoming request body
	var signatureReq model.SignatureData

	// Parse the request body
	if err := c.BodyParser(&signatureReq); err != nil {
		return BadRequest(c, "Gagal memparsing body request! Silakan periksa kembali.", err.Error())
	}

	// Validate the input data
	if _, err := govalidator.ValidateStruct(signatureReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Silakan periksa kembali.", "Data tidak valid")
	}

	// Retrieve the admin ID from the claims stored in context
	claims := c.Locals("admin").(jwt.MapClaims)
	adminID, ok := claims["sub"].(string)
	if !ok {
		return Unauthorized(c, "Token Admin tidak valid!", "Token Admin tidak valid!")
	}

	// Convert adminID (which is a string) to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return Unauthorized(c, "Format token admin tidak valid!", "Format token admin tidak valid!")
	}

	// Create a new signature object
	signature := model.Signature{
		SignatureData: model.SignatureData{
			AdminId:    objectID,
			ConfigName: signatureReq.ConfigName,
			Stamp:      signatureReq.Stamp,
			Signature:  signatureReq.Signature,
			Logo:       signatureReq.Logo,
			Name:       signatureReq.Name,
			Role:       signatureReq.Role,
		},
		Model: model.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// Insert the new signature into the database
	_, err = collectionSignature.InsertOne(context.TODO(), signature)
	if err != nil {
		return Conflict(c, "Gagal membuat signature baru! Silakan coba lagi.", "Gagal menambah")
	}

	// Return success response
	return OK(c, "Berhasil membuat signature baru!", signature)
}

// Function to get all signatures
func GetAllSignature(c *fiber.Ctx) error {
	var results []bson.M // Slice to hold the results
	ctx := c.Context()

	// Retrieve the admin ID from the claims stored in context
	claims := c.Locals("admin").(jwt.MapClaims)
	adminID, ok := claims["sub"].(string)
	if !ok {
		return Unauthorized(c, "Token Admin tidak valid!", "Token Admin tidak valid!")
	}

	// Convert adminID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return Unauthorized(c, "Format token admin tidak valid!", "Format token admin tidak valid!")
	}

	// Set the projection to return the required fields
	projection := bson.M{
		"_id":         1,
		"admin_id":    1,
		"config_name": 1,
		"created_at":  1,
		"updated_at":  1,
		"deleted_at":  1,
	}

	// Create the filter to include admin_id and handle deleted_at
	filter := bson.M{
		"admin_id": objectID,
		"$or": []bson.M{
			{"deleted_at": bson.M{"$exists": false}}, // DeletedAt field does not exist
			{"deleted_at": bson.M{"$eq": nil}},       // DeletedAt field is nil
		},
	}

	// Find all documents that match
	cursor, err := collectionSignature.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Tidak ada signature yang ditemukan!", err.Error())
		}
		return Conflict(c, "Gagal mengambil data signature! Silakan coba lagi.", err.Error())
	}
	defer cursor.Close(ctx) // Close the cursor after use

	// Decode each document and append it to results
	for cursor.Next(ctx) {
		var signature bson.M
		if err := cursor.Decode(&signature); err != nil {
			return Conflict(c, "Gagal mengambil data signature! Silakan coba lagi.", "")
		}
		results = append(results, signature) // Append the signature to results
	}
	if err := cursor.Err(); err != nil {
		return Conflict(c, "Gagal menampilkan semua data signature! Silakan coba lagi.", "")
	}

	// Return success response with all signatures
	return OK(c, "Berhasil menampilkan semua data signature!", results)
}

// Function to get a single signature
func GetSignatureByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	searchKey := c.Params("type")
	if searchKey == "" { // from handler w/o type param, to not break api
		searchKey = "oid"
	}
	var searchVal any
	searchVal = idParam

	// Convert to ObjectID if needed
	if searchKey == "oid" {
		searchKey = "_id"
		certifID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return BadRequest(c, "Signature ini tidak ada!", "Please provide a valid ObjectID")
		}
		searchVal = certifID
	}

	// Make filter to find document based on search key & value
	filter := bson.M{searchKey: searchVal}

	// Variable to hold the signature details
	var signatureDetail bson.M

	// Find a single document that matches the filter
	if err := collectionSignature.FindOne(context.TODO(), filter).Decode(&signatureDetail); err != nil {
		if err == mongo.ErrNoDocuments {
			// If not found, return a 404 status
			return NotFound(c, "Signature ini tidak dapat ditemukan! Silakan periksa ID yang dimasukkan.", err.Error())
		}
		return Conflict(c, "Gagal mendapatkan data! Silakan coba lagi.", "KONFLIK!")
	}

	// Check if the signature has been deleted
	if deletedAt, exists := signatureDetail["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Signature ini telah dihapus! Silakan hubungi admin.", "Periksa signature yang dihapus", deletedAt)
	}

	// Return success response with signature details
	return OK(c, "Berhasil menampilkan data signature!", signatureDetail)
}

// Function to edit a signature
func EditSignature(c *fiber.Ctx) error {
	idParam := c.Params("id") // Get ID from the URL parameters

	// Convert ID to ObjectID
	signatureID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Signature ini tidak ada! Silakan periksa ID yang dimasukkan.", "")
	}

	filter := bson.M{"_id": signatureID} // Create filter to find the signature

	var input model.SignatureData

	var signatureData bson.M // Variable to hold the found signature data

	// Find the signature based on the filter
	if err := collectionSignature.FindOne(c.Context(), filter).Decode(&signatureData); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Signature ini tidak dapat ditemukan! Silakan periksa ID yang dimasukkan.", "")
		}
		return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", "")
	}

	// Check if the signature has been deleted
	if deletedAt, exists := signatureData["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Signature ini telah dihapus! Silakan hubungi admin.", "Periksa signature yang dihapus", deletedAt)
	}

	// Parse the request body to get new data
	if err := c.BodyParser(&input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Silakan periksa kembali.", err.Error())
	}

	// Validate the input data
	if _, err := govalidator.ValidateStruct(input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Silakan periksa kembali.", "Data tidak valid")
	}

	// Create update for the signature data
	update := bson.M{
		"config_name": input.ConfigName,
		"stamp":       input.Stamp,
		"signature":   input.Signature,
		"name":        input.Name,
		"role":        input.Role,
		"updated_at":  time.Now(),
	}

	// Update the signature in the database
	_, err = collectionSignature.UpdateOne(c.Context(), filter, bson.M{"$set": update})
	if err != nil {
		return Conflict(c, "Gagal memperbarui signature! Silakan coba lagi.", err.Error())
	}

	// Return success response
	return OK(c, "Berhasil memperbarui signature!", update)
}

// Function to delete a signature
func DeleteSignature(c *fiber.Ctx) error {
	idParam := c.Params("id") // Get ID from the URL parameters

	// Convert ID to ObjectID
	signatureID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Signature ini tidak ada! Silakan periksa ID yang dimasukkan.", "Gagal mengonversi ID")
	}

	filter := bson.M{"_id": signatureID} // Create filter to find the signature

	var signatureData bson.M // Variable to hold the found signature data
	err = collectionSignature.FindOne(context.TODO(), filter).Decode(&signatureData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Tidak dapat menemukan signature! Silakan periksa ID yang dimasukkan.", "Gagal menemukan signature")
		}
		return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", "Gagal menemukan signature")
	}

	// Check if the signature has been deleted
	if deletedAt, exists := signatureData["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Signature ini telah dihapus! Silakan hubungi admin.", "Periksa signature yang dihapus", deletedAt)
	}

	// Update the deleted_at timestamp
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	// Update the signature in the database
	result, err := collectionSignature.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal menghapus signature! Silakan coba lagi.", "Gagal menghapus signature")
	}

	// Check if the document was found and updated
	if result.MatchedCount == 0 {
		return NotFound(c, "Signature ini tidak dapat ditemukan! Silakan periksa ID yang dimasukkan.", "Gagal menemukan signature")
	}

	// Return success response
	return OK(c, "Berhasil menghapus signature!", signatureID)
}
