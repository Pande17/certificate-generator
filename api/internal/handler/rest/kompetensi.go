package rest

import (
	"context"
	"time"

	"certificate-generator/database"
	"certificate-generator/model"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to the competence collection in the database
var collectionKompetensi = database.GetCollection("competence")

// Function to create a new competence
func CreateKompetensi(c *fiber.Ctx) error {
	// Struct for the incoming request body
	var kompetensiReq struct {
		KompetensiName string        `json:"nama_kompetensi" valid:"required~Nama Kompetensi tidak boleh kosong!, stringlength(3|50)~Nama Kompetensi harus antara 3-50 karakter"`
		Divisi         string        `json:"divisi" valid:"required~Divisi tidak boleh kosong!, stringlength(1|6)~Divisi harus antara 1-6 karakter"`
		SKKNI          string        `json:"skkni" valid:"required~SKKNI tidak boleh kosong!"`
		HardSkills     []model.Skill `json:"hard_skills"`
		SoftSkills     []model.Skill `json:"soft_skills"`
	}

	// Parse the request body
	if err := c.BodyParser(&kompetensiReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Mohon periksa kembali.", "Data yang dimasukkan tidak valid!")
	}

	// Validate the input data using govalidator
	if _, err := govalidator.ValidateStruct(&kompetensiReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", err.Error())
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

	// Check the availability of the competence name
	var existingKompetensi model.Kompetensi
	filter := bson.M{"nama_kompetensi": kompetensiReq.KompetensiName}

	// Find competence with the same name
	err = collectionKompetensi.FindOne(context.TODO(), filter).Decode(&existingKompetensi)
	if err == nil {
		return Conflict(c, "Kompetensi dengan nama ini sudah ada! Silakan gunakan nama lain.", "Kompetensi dengan nama yang sama sudah ada!")
	} else if err != mongo.ErrNoDocuments {
		return Conflict(c, "Gagal memeriksa kompetensi yang ada. Silakan coba lagi.", err.Error())
	}

	// Create a new competence
	kompetensi := model.Kompetensi{
		AdminId:        objectID,
		NamaKompetensi: kompetensiReq.KompetensiName,
		Divisi:         kompetensiReq.Divisi,
		SKKNI:          kompetensiReq.SKKNI,
		HardSkills:     kompetensiReq.HardSkills,
		SoftSkills:     kompetensiReq.SoftSkills,
		Model: model.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// Insert the new competence into the database
	_, err = collectionKompetensi.InsertOne(context.TODO(), kompetensi)
	if err != nil {
		return Conflict(c, "Gagal membuat data kompetensi baru! Silakan coba lagi.", "Gagal menyimpan kompetensi baru")
	}

	// Return success
	return OK(c, "Berhasil membuat Kompetensi Baru!", kompetensi)
}

// Function to edit competence
func EditKompetensi(c *fiber.Ctx) error {
	// Get kompetensi_id from params
	idParam := c.Params("id")

	// Convert kompetensi_id to ObjectID
	kompetensiID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Kompetensi ini tidak ada! Silakan periksa ID yang dimasukkan.", "Gagal mengonversi ID")
	}

	// Create filter to find the document
	filter := bson.M{"_id": kompetensiID}

	// Variable to hold results
	var competenceData bson.M

	// Search for the competence based on its ID
	if err := collectionKompetensi.FindOne(c.Context(), filter).Decode(&competenceData); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Kompetensi ini tidak dapat ditemukan! Silakan periksa ID yang dimasukkan.", "Gagal menemukan kompetensi")
		}
		return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", "Gagal mengambil kompetensi")
	}

	// Check if the competence has been deleted
	if deletedAt, exists := competenceData["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Kompetensi ini telah dihapus! Silakan hubungi admin.", "Periksa kompetensi yang dihapus", deletedAt)
	}

	// Parse request body to get new data
	var input struct {
		NamaKompetensi string        `json:"nama_kompetensi" valid:"required~Nama Kompetensi tidak boleh kosong!, stringlength(3|50)~Nama Kompetensi harus antara 3-50 karakter"`
		Divisi         string        `json:"divisi" valid:"required~Divisi tidak boleh kosong!, stringlength(1|6)~Divisi harus antara 1-6 karakter"`
		SKKNI          string        `json:"skkni" valid:"required~SKKNI tidak boleh kosong!"`
		HardSkills     []model.Skill `json:"hard_skills"`
		SoftSkills     []model.Skill `json:"soft_skills"`
	}

	// Handle if request body is invalid
	if err := c.BodyParser(&input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Silakan periksa kembali.", "Gagal mem-parsing body")
	}

	// Validate the input data using govalidator
	if _, err := govalidator.ValidateStruct(&input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", err.Error())
	}

	// Update fields in the database
	update := bson.M{
		"$set": bson.M{
			"nama_kompetensi": input.NamaKompetensi,
			"divisi":          input.Divisi,
			"skkni":           input.SKKNI,
			"hard_skills":     input.HardSkills,
			"soft_skills":     input.SoftSkills,
			"updated_at":      time.Now(),
		},
	}

	// Update data in the collection
	_, err = collectionKompetensi.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal memperbarui Kompetensi! Silakan coba lagi.", "Gagal memperbarui kompetensi")
	}

	// Return success
	return OK(c, "Berhasil memperbarui Kompetensi!", update)
}

// Function to delete competence
func DeleteKompetensi(c *fiber.Ctx) error {
	// Get id param
	idParam := c.Params("id")

	// Convert idParam to ObjectID
	kompetensiID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Kompetensi ini tidak ada! Silakan periksa ID yang dimasukkan.", "Gagal mengonversi ID")
	}

	// Create filter to find the document
	filter := bson.M{"_id": kompetensiID}

	// Find competence
	var competenceData bson.M
	err = collectionKompetensi.FindOne(context.TODO(), filter).Decode(&competenceData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Tidak dapat menemukan Kompetensi! Silakan periksa ID yang dimasukkan.", "Gagal menemukan kompetensi")
		}
		return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", "Gagal mengambil kompetensi")
	}

	// Check if the competence has been deleted
	if deletedAt, exists := competenceData["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Kompetensi ini telah dihapus! Silakan hubungi admin.", "Periksa kompetensi yang dihapus", deletedAt)
	}

	// Update the deleted_at timestamp
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	// Update document in the collection
	result, err := collectionKompetensi.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal menghapus kompetensi! Silakan coba lagi.", "Gagal menghapus kompetensi")
	}

	// Check if the document was already deleted
	if result.MatchedCount == 0 {
		return NotFound(c, "Kompetensi ini tidak dapat ditemukan! Silakan periksa ID yang dimasukkan.", "Gagal menemukan kompetensi")
	}

	// Return success
	return OK(c, "Berhasil menghapus kompetensi!", kompetensiID)
}

// Function to get all competencies
func GetAllKompetensi(c *fiber.Ctx) error {
	var results []bson.M
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
		"_id":             1,
		"admin_id":        1,
		"nama_kompetensi": 1,
		"divisi":          1,
		"skkni":           1,
		"created_at":      1,
		"updated_at":      1,
		"deleted_at":      1,
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
	cursor, err := collectionKompetensi.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Kompetensi tidak dapat ditemukan!", err.Error())
		}
		return Conflict(c, "Gagal mengambil data kompetensi! Silakan coba lagi.", err.Error())
	}
	defer cursor.Close(ctx)

	// Decode each document and append it to results
	for cursor.Next(ctx) {
		var competence bson.M
		if err := cursor.Decode(&competence); err != nil {
			return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", "Gagal mendekode kompetensi")
		}
		results = append(results, competence)
	}
	if err := cursor.Err(); err != nil {
		return Conflict(c, "Gagal menampilkan data! Silakan coba lagi.", "Gagal menampilkan kompetensi")
	}

	// Return success
	return OK(c, "Berhasil menampilkan semua data Kompetensi!", results)
}

// Function to get one competence by filter
func GetKompetensiByID(c *fiber.Ctx) error {
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
			return BadRequest(c, "Sertifikat ini tidak ada!", "Please provide a valid ObjectID")
		}
		searchVal = certifID
	}

	// Make filter to find document based on search key & value
	filter := bson.M{searchKey: searchVal}

	// Variable to hold search results
	var kompetensiDetail bson.M

	// Find a single document that matches the filter
	if err := collectionKompetensi.FindOne(context.TODO(), filter).Decode(&kompetensiDetail); err != nil {
		// If not found, return a 404 status
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Kompetensi ini tidak dapat ditemukan! Silakan periksa ID yang dimasukkan.", "Gagal menemukan detail kompetensi")
		}
		return Conflict(c, "Gagal mendapatkan data! Silakan coba lagi.", "Gagal menemukan detail kompetensi")
	}

	// Check if the competence has been deleted
	if deletedAt, exists := kompetensiDetail["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Kompetensi ini telah dihapus! Silakan hubungi admin.", "Periksa kompetensi yang dihapus", deletedAt)
	}

	// Return success
	return OK(c, "Berhasil menampilkan data kompetensi!", kompetensiDetail)
}
