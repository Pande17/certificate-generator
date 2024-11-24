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

var collectionSignature = database.GetCollection("signature")

func CreateSignature(c *fiber.Ctx) error {
	var signatureReq struct {
		ConfigName string `json:"config_name" bson:"config_name" valid:"required~Stamp tidak boleh kosong!"`
		Stamp      string `json:"stamp" valid:"required~Stamp tidak boleh kosong!, url"`
		Signature  string `json:"signature" valid:"required~Signature tidak boleh kosong!, url"`
		Name       string `json:"name" valid:"required~Nama tidak boleh kosong!, stringlength(1|60)~Nama harus antara 1 hingga 60 karakter!"`
		Role       string `json:"role" valid:"required~Role tidak boleh kosong!, stringlength(1|60)~Role harus antara 1 hingga 60 karakter!"`
	}

	if err := c.BodyParser(&signatureReq); err != nil {
		return BadRequest(c, "Tidak dapat memparsing body request!", err.Error())
	}

	if _, err := govalidator.ValidateStruct(signatureReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", "Data tidak valid")
	}

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

	_, err := collectionSignature.InsertOne(context.TODO(), signature)
	if err != nil {
		return Conflict(c, "Gagal membuat Signature, silahkan coba lagi!", "")
	}

	return OK(c, "Berhasil membuat Signature baru!", signature)
}

func GetSignature(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return getAllSignature(c)
	}

	var value any
	var err error
	value, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return BadRequest(c, "Gagal mendapatkan Signature!", "")
	}
	return getOneSignature(c, bson.M{"_id": value})
}

func getAllSignature(c *fiber.Ctx) error {
	var results []bson.M

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projection := bson.M{
		"_id":         1,
		"config_name": 1,
		"created_at":  1,
		"updated_at":  1,
		"deleted_at":  1,
	}

	cursor, err := collectionSignature.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Tidak ada Signature yang dibuat!", err.Error())
		}
		return Conflict(c, "Gagal mengambil data Signature!", err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var signature bson.M
		if err := cursor.Decode(&signature); err != nil {
			return Conflict(c, "Gagal mengambil data Signature!", "")
		}
		if deletedAt, ok := signature["deleted_at"]; ok && deletedAt != nil {
			continue
		}
		results = append(results, signature)
	}
	if err := cursor.Err(); err != nil {
		return Conflict(c, "Gagal menampilkan semua data Signature!", "")
	}

	return OK(c, "Berhasil menampilkan semua data Signature!", results)
}

func getOneSignature(c *fiber.Ctx, filter bson.M) error {
	var signatureDetail bson.M

	if err := collectionSignature.FindOne(context.TODO(), filter).Decode(&signatureDetail); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Signature ini tidak dapat ditemukan!", err.Error())
		}
		return Conflict(c, "Gagal mendapatkan data!", "KONFLIK!")
	}

	if deletedAt, exists := signatureDetail["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Signature ini telah dihapus", "", deletedAt)
	}

	return OK(c, "Berhasil menampilkan data Signature!", signatureDetail)
}

func EditSignature(c *fiber.Ctx) error {
	idParam := c.Params("id")

	signatureID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Signature ini tidak ada!", "")
	}

	filter := bson.M{"_id": signatureID}

	var input struct {
		ConfigName string `json:"config_name" bson:"config_name" valid:"required~Stamp tidak boleh kosong!"`
		Stamp      string `json:"stamp" valid:"required~Stamp tidak boleh kosong!, url"`
		Signature  string `json:"signature" valid:"required~Signature tidak boleh kosong!, url"`
		Name       string `json:"name" valid:"required~Nama tidak boleh kosong!, stringlength(1|60)~Nama harus antara 1 hingga 60 karakter!"`
		Role       string `json:"role" valid:"required~Role tidak boleh kosong!, stringlength(1|60)~Role harus antara 1 hingga 60 karakter!"`
	}

	var signatureData bson.M

	if err := collectionSignature.FindOne(c.Context(), filter).Decode(&signatureData); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Signature ini tidak dapat ditemukan!", "")
		}
		return Conflict(c, "Gagal mengambil Data!", "")
	}

	if deletedAt, exists := signatureData["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Signature ini telah dihapus!", "Check deleted Signature", deletedAt)
	}

	if err := c.BodyParser(&input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", err.Error())
	}

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

	_, err = collectionSignature.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal memperbarui Signature!", err.Error())
	}

	return OK(c, "Berhasil memperbarui Signature!", update)
}

func DeleteSignature(c *fiber.Ctx) error {
	idParam := c.Params("id")

	signatureID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Signature ini tidak ada!", "Convert Params Delete Signature")
	}

	filter := bson.M{"_id": signatureID}

	var signatureData bson.M
	err = collectionSignature.FindOne(context.TODO(), filter).Decode(&signatureData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Tidak dapat menemukan Signature!", "Find Signature")
		}
		return Conflict(c, "Gagal mengambil data!", "Find Signature")
	}

	if deletedAt, exists := signatureData["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Signature ini telah dihapus!", "Check deleted Signature", deletedAt)
	}

	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	result, err := collectionSignature.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal menghapus Signature!", "Delete Signature")
	}

	if result.MatchedCount == 0 {
		return NotFound(c, "Signature ini tidak dapat ditemukan!", "Check deleted Signature on Delete")
	}

	return OK(c, "Berhasil menghapus Signature!", signatureID)
}
