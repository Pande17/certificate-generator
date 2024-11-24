package rest

import (
	"context"
	"fmt"

	"certificate-generator/database"
	"certificate-generator/model"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect collection competence in database
var collectionKompetensi = database.GetCollection("competence")

// func for create new kompetensi
func CreateKompetensi(c *fiber.Ctx) error {
	// struct for the incoming request body
	var kompetensiReq struct {
		KompetensiName string        `json:"nama_kompetensi" valid:"required~Nama Kompetensi tidak boleh kosong!, stringLength(3|50)~Nama Kompetensi harus antara 3-50 karakter"`
		Divisi         string        `json:"divisi" valid:"required~Divisi tidak boleh kosong!, stringLength(1|6)~Divisi harus antara 1-6 karakter"`
		HardSkills     []model.Skill `json:"hard_skills"`
		SoftSkills     []model.Skill `json:"soft_skills"`
	}

	// parse the request body
	if err := c.BodyParser(&kompetensiReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Mohon diperiksa kembali!", "Data yang dimasukkan tidak valid!")
	}

	// Validate the input data using govalidator
	if _, err := govalidator.ValidateStruct(&kompetensiReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", err.Error())
	}

	// Retrieve the user ID from the claims stored in context
	claims := c.Locals("admin").(jwt.MapClaims)
	adminID, ok := claims["sub"].(string)
	if !ok {
		return Unauthorized(c, "Token Admin tidak valid!", "Token Admin tidak valid!")
	}

	// Convert userID (which is a string) to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(adminID)
	fmt.Println("Admin ID from token:", adminID)
	if err != nil {
		return Unauthorized(c, "Format token admin tidak valid!", "Formaat token admin tidak valid!")
	}

	// new variable to check the availability of the competence name
	var existingKompetensi model.Kompetensi

	// new variable to find competence based on their name "competence_name"
	filter := bson.M{"nama_kompetensi": kompetensiReq.KompetensiName}

	// find competence with same competence name as input name
	err = collectionKompetensi.FindOne(context.TODO(), filter).Decode(&existingKompetensi)
	if err == nil {
		return Conflict(c, "Kompetensi dengan nama ini sudah ada!", "Kompetensi dengan nama yang sama sudah ada!")
	} else if err != mongo.ErrNoDocuments {
		return Conflict(c, "Gagal dalam memeriksa Kompetensi yang ada", err.Error())
	}

	// append data from body request to struct Kompetensi
	kompetensi := model.Kompetensi{
		AdminId:        objectID,
		NamaKompetensi: kompetensiReq.KompetensiName,
		Divisi:         kompetensiReq.Divisi,
		HardSkills:     kompetensiReq.HardSkills,
		SoftSkills:     kompetensiReq.SoftSkills,
		Model: model.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// insert data from struct "Kompetensi" to collection "competence" in database MongoDB
	_, err = collectionKompetensi.InsertOne(context.TODO(), kompetensi)
	if err != nil {
		return Conflict(c, "Gagal membuat data kompetensi yang baru!", "Failed input new competence")
	}

	// return success
	return OK(c, "Berhasil membuat Kompetensi Baru!", kompetensi)
}

func EditKompetensi(c *fiber.Ctx) error {
	// get kompetensi_id from params
	idParam := c.Params("id")

	// convert kompetensi_id to integer data type
	kompetensiID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Kompetensi ini tidak ada!", "Convert Params")
	}

	// connect to collection in MongoDB
	collectionKompetensi := database.GetCollection("competence")

	// make filter to find document based on params
	filter := bson.M{"_id": kompetensiID}

	// variabwle to hold results
	var competenceData bson.M

	// searching for the competence based on their kompetensi_id
	if err := collectionKompetensi.FindOne(c.Context(), filter).Decode(&competenceData); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Kompetensi ini tidak dapat ditemukan!", "Find kompetensi_id based on params")
		}
		return Conflict(c, "Gagal mengambil data!", "Find kompetensi_id based on params")
	}

	// Modified code for DeleteKompetensi
	// if modelData, ok := competenceData["model"].(bson.M); ok {
	if deletedAt, exists := competenceData["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Kompetensi ini telah dihapus!", "Check deleted kompetensi", deletedAt)
	}
	// }

	// parsing req body to get new data
	var input struct {
		NamaKompetensi string        `json:"nama_kompetensi"`
		Divisi         string        `json:"divisi"`
		HardSkills     []model.Skill `json:"hard_skills"`
		SoftSkills     []model.Skill `json:"soft_skills"`
	}

	// handler if request body is invalid
	if err := c.BodyParser(&input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", "Req body edit Kompetensi")
	}

	// update fields in the database
	update := bson.M{
		"$set": bson.M{
			"nama_kompetensi": input.NamaKompetensi,
			"divisi":          input.Divisi,
			"hard_skills":     input.HardSkills,
			"soft_skills":     input.SoftSkills,
			"updated_at":      time.Now(),
		},
	}

	// update data in collection based on their "kompetensi_id" or params
	_, err = collectionKompetensi.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal memperbarui Kompetensi!", "Update new data kompetensi")
	}

	// return success
	return OK(c, "Berhasil memperbarui Kompetensi!", update)
}

func DeleteKompetensi(c *fiber.Ctx) error {
	// get kompetensi_id from params
	idParam := c.Params("id")

	// convert params to integer data type
	kompetensiID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Kompetensi ini tidak ada!", "Convert Params Delete Kompetensi")
	}

	// connect to collection in MongoDB
	collectionKompetensi := database.GetCollection("competence")

	// make filter to find document based on kompetensi_id
	filter := bson.M{"_id": kompetensiID}

	// find competence
	var competenceData bson.M
	err = collectionKompetensi.FindOne(context.TODO(), filter).Decode(&competenceData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Tidak dapat menemukan Kompetensi!", "Find Kompetensi")
		}
		return Conflict(c, "Gagal mengambil data!", "Find Kompetensi")
	}

	// Modified code for DeleteKompetensi
	// if modelData, ok := competenceData["model"].(bson.M); ok {
	if deletedAt, exists := competenceData["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Kompetensi ini telah dihapus!", "Check deleted kompetensi", deletedAt)
	}
	// }

	// make update for input timestamp DeletedAt
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	// update document in collection MongoDB
	result, err := collectionKompetensi.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal menghapus Kompetensi!", "Delete Kompetensi")
	}

	// check if the document is found and updated
	if result.MatchedCount == 0 {
		return NotFound(c, "Kompetensi ini tidak dapat ditemukan!", "Check deleted kompetensi on Delete")
	}

	// return success
	return OK(c, "Berhasil menghapus Kompetensi!", kompetensiID)
}

// function to get all kompetensi data
func GetKompetensi(c *fiber.Ctx) error {
	id := c.Params("id") // Get ID from the URL path
	if id == "" {
		// If ID is not provided, return all kompetensi data
		return getAllKompetensi(c)
	}

	// If ID is provided, proceed with getting specific kompetensi
	var value any
	var err error
	value, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return BadRequest(c, "Gagal mendapatkan Kompetensi!", err.Error())
	}
	return getOneKompetensi(c, bson.M{"_id": value})
}

func getAllKompetensi(c *fiber.Ctx) error {
	var results []bson.M

	collectionKompetensi := database.GetCollection("competence")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set the projection to return the required fields
	projection := bson.M{
		"_id":             1, // 0 to exclude the field
		"admin_id":        1,
		"nama_kompetensi": 1, // 1 to include the field, _id will be included by default
		"divisi":          1,
		"created_at":      1,
		"updated_at":      1,
		"deleted_at":      1,
	}

	// find the projection
	cursor, err := collectionKompetensi.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Kompetensi tidak dapat ditemukan!", err.Error())
		}
		return Conflict(c, "Gagal mengambil data kompetensi!", err.Error())
	}
	defer cursor.Close(ctx)

	// decode each document and append it to results
	for cursor.Next(ctx) {
		var competence bson.M
		if err := cursor.Decode(&competence); err != nil {
			return Conflict(c, "Gagal mengambil data", "Decode Kompetensi")
		}
		if deletedAt, ok := competence["deleted_at"]; ok && deletedAt != nil {
			// skip deleted certificates
			continue
		}
		results = append(results, competence)
	}
	if err := cursor.Err(); err != nil {
		return Conflict(c, "Gagal menampilkan data!", "Append Kompetensi")
	}

	// return success
	return OK(c, "Berhasil menampilkan semua data Kompetensi!", results)
}

func getOneKompetensi(c *fiber.Ctx, filter bson.M) error {
	// connect to collection in MongoDB
	collectionKompetensi := database.GetCollection("competence")

	// variable to hold search results
	var kompetensiDetail bson.M

	// find a single document that matches the filter
	if err := collectionKompetensi.FindOne(context.TODO(), filter).Decode(&kompetensiDetail); err != nil {
		// if not found, return a 404 status
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Kompetensi ini tidak dapat ditemukan!", "Find Detail Kompetensi")
		}
		// if in server error, return status 500
		return Conflict(c, "Gagal mendapatkan data!", "Server Find Detail Kompetensi")
	}

	// Check if the competence has a "deleted_at" field
	if deletedAt, exists := kompetensiDetail["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Kompetensi ini telah dihapus!", "Check deleted kompetensi on get Detail", deletedAt)
	}

	// return success
	return OK(c, "Berhasil menampilkan data Kompetensi!", kompetensiDetail)
}
