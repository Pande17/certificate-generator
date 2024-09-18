package controller

import (
	"context"
	"pkl/finalProject/certificate-generator/repository/config"
	dbmongo "pkl/finalProject/certificate-generator/repository/db_mongo"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// func for create new kompetensi
func CreateKompetensi(c *fiber.Ctx) error {
	// struct for the incoming request body
	var kompetensiReq struct {
		ID             primitive.ObjectID  `bson:"_id,omitempty"`
		KompetensiID   uint64              `bson:"kompetensi_id"`
		KompetensiName string              `json:"nama_kompetensi"`
		HardSkills     []dbmongo.HardSkill `json:"hard_skills"`
		SoftSkills     []dbmongo.SoftSkill `json:"soft_skills"`
	}

	// parse the request body
	if err := c.BodyParser(&kompetensiReq); err != nil {
		return BadRequest(c, "Failed to read body")
	}

	// connect collection competence in database
	collection := config.MongoClient.Database("certificate-generator").Collection("competence")

	// new variable to check the availability of the competence name
	var existingKompetensi dbmongo.Kompetensi

	// new variable to find competence based on their name "competence_name"
	filter := bson.M{"nama_kompetensi": kompetensiReq.KompetensiName}

	// find competence with same competence name as input name
	err := collection.FindOne(context.TODO(), filter).Decode(&existingKompetensi)
	if err == nil {
		return Conflict(c, "Competence already exists")
	} else if err != mongo.ErrNoDocuments {
		return InternalServerError(c, "Error chechking for existing Competence")
	}

	// generate kompetensi_id (incremental id)
	nextKompetensiID, err := GetNextIncrementalID(collection, "kompetensi_id")
	if err != nil {
		return InternalServerError(c, "Failed to generate Kompetensi ID")
	}

	//
	kompetensi := dbmongo.Kompetensi{
		ID:             primitive.NewObjectID(),
		KompetensiID:   uint64(nextKompetensiID),
		NamaKompetensi: kompetensiReq.KompetensiName,
		HardSkills:     kompetensiReq.HardSkills,
		SoftSkills:     kompetensiReq.SoftSkills,
		Model: dbmongo.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// insert data from struct "Kompetensi" to collection "competence" in database MongoDB
	_, err = collection.InsertOne(context.TODO(), kompetensi)
	if err != nil {
		return InternalServerError(c, "Failed to create New Competence")
	}

	// return success
	return Ok(c, "Sucess created New Competence", kompetensi)
}

func EditKompetensi(c *fiber.Ctx) error {

	// return success
	return Ok(c, "Sucess edited Competence data", nil)
}

func DeleteKompetensi(c *fiber.Ctx) error {

	// return success
	return Ok(c, "Sucess deleted Competence data", nil)
}

func SeeAllKompetensi(c *fiber.Ctx) error {

	// return success
	return Ok(c, "Sucess get all Competence data", nil)
}

func SeeDetailKompetensi(c *fiber.Ctx) error {

	// return success
	return Ok(c, "Sucess get Competence data", nil)
}
