package rest

import (
	"certificate-generator/database"
	"certificate-generator/model"
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
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
	var signatureReq struct {
		ConfigName string `json:"config_name" bson:"config_name" valid:"required~Nama konfigurasi tidak boleh kosong!"`
		Stamp      string `json:"stamp" valid:"required~Stamp tidak boleh kosong!, url"`
		Signature  string `json:"signature" valid:"required~Signature tidak boleh kosong!, url"`
		Name       string `json:"name" valid:"required~Nama tidak boleh kosong!, stringlength(1|60)~Nama harus antara 1 hingga 60 karakter!"`
		Role       string `json:"role" valid:"required~Peran tidak boleh kosong!, stringlength(1|60)~Peran harus antara 1 hingga 60 karakter!"`
	}

	// Parse the request body
	if err := c.BodyParser(&signatureReq); err != nil {
		return BadRequest(c, "Gagal memparsing body request! Silakan periksa kembali.", err.Error())
	}

	// Validate the input data
	if _, err := govalidator.ValidateStruct(signatureReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Silakan periksa kembali.", "Data tidak valid")
	}

	// Create a new signature object
	signature := model.Signature{
		ConfigName: signatureReq.ConfigName,
		Stamp:      signatureReq.Stamp,
		Signature:  signatureReq.Signature,
		Name:       signatureReq.Name,
		Role:       signatureReq.Role,
		Model: model.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// Insert the new signature into the database
	_, err := collectionSignature.InsertOne(context.TODO(), signature)
	if err != nil {
		return Conflict(c, "Gagal membuat signature baru! Silakan coba lagi.", "")
	}

	// Return success response
	return OK(c, "Berhasil membuat signature baru!", signature)
}

// Function to get a signature by ID or all signatures
func GetSignature(c *fiber.Ctx) error {
	id := c.Params("id") // Get ID from the URL parameters
	if id == "" {
		// If no ID is provided, return all signatures
		return getAllSignature(c)
	}

	// If ID is provided, proceed to get the specific signature
	var value any
	var err error
	value, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return BadRequest(c, "Gagal mendapatkan signature! Silakan periksa ID yang dimasukkan.", "")
	}
	return getOneSignature(c, bson.M{"_id": value})
}

// Function to get all signatures
func getAllSignature(c *fiber.Ctx) error {
	var results []bson.M // Slice to hold the results

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Cancel the context after the function completes

	// Set the projection to return the required fields
	projection := bson.M{
		"_id":         1,
		"config_name": 1,
		"created_at":  1,
		"updated_at":  1,
		"deleted_at":  1,
	}

	// Find all signatures in the collection
	cursor, err := collectionSignature.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Tidak ada signature yang ditemukan!", err.Error())
		}
		return Conflict(c, "Gagal mengambil data signature! Silakan coba lagi.", err.Error())
	}
	defer cursor.Close(ctx) // Close the cursor after use

	// Iterate through each result in the cursor
	for cursor.Next(ctx) {
		var signature bson.M
		if err := cursor.Decode(&signature); err != nil {
			return Conflict(c, "Gagal mengambil data signature! Silakan coba lagi.", "")
		}
		// Skip signatures that have been deleted
		if deletedAt, ok := signature["deleted_at"]; ok && deletedAt != nil {
			continue
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
func getOneSignature(c *fiber.Ctx, filter bson.M) error {
	var signatureDetail bson.M // Variable to hold the signature details

	// Find a single document that matches the filter
	if err := collectionSignature.FindOne(context.TODO(), filter).Decode(&signatureDetail); err != nil {
		if err == mongo.ErrNoDocuments {
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

	var input struct {
		ConfigName string `json:"config_name" bson:"config_name" valid:"required~Nama konfigurasi tidak boleh kosong!"`
		Stamp      string `json:"stamp" valid:"required~Stamp tidak boleh kosong!, url"`
		Signature  string `json:"signature" valid:"required~Signature tidak boleh kosong!, url"`
		Name       string `json:"name" valid:"required~Nama tidak boleh kosong!, stringlength(1|60)~Nama harus antara 1 hingga 60 karakter!"`
		Role       string `json:"role" valid:"required~Peran tidak boleh kosong!, stringlength(1|60)~Peran harus antara 1 hingga 60 karakter!"`
	}

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
		"$set": bson.M{
			"config_name": input.ConfigName,
			"stamp":       input.Stamp,
			"signature":   input.Signature,
			"name":        input.Name,
			"role":        input.Role,
			"updated_at":  time.Now(),
		},
	}

	// Update the signature in the database
	_, err = collectionSignature.UpdateOne(c.Context(), filter, update)
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
