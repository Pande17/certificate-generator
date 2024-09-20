package controller

import (
	"context"
	"fmt"

	"pkl/finalProject/certificate-generator/repository/config"
	dbmongo "pkl/finalProject/certificate-generator/repository/db_mongo"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// get kompetensi_id from params
	idParam := c.Params("id")

	// convert kompetensi_id to integer data type
	kompetensiID, err := strconv.Atoi(idParam)
	if err != nil {
		return BadRequest(c, "Invalid ID")
	}

	// connect to collection in MongoDB
	collection := config.MongoClient.Database("certificate-generator").Collection("competence")

	// make filter to find document based on params
	filter := bson.M{"kompetensi_id": kompetensiID}

	// variabwle to hold results
	var competenceData bson.M

	// searching for the competence based on their kompetensi_id
	if err := collection.FindOne(c.Context(), filter).Decode(&competenceData); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Competence not found")
		}
		return InternalServerError(c, "Failed to fetch data")
	}

	// check if competence has already been deleted
	if deletedAt, ok := competenceData["model"].(bson.M)["deleted_at"]; ok && deletedAt != nil {
		// return the deletion time if the competence is already deleted
		return AlreadyDeleted(c, "This competence has already been Deleted", deletedAt)
	}

	// parsing req body to get new data
	var input struct {
		NamaKompetensi string              `json:"nama_kompetensi"`
		HardSkills     []dbmongo.HardSkill `json:"hard_skills"`
		SoftSkills     []dbmongo.SoftSkill `json:"soft_skills"`
	}

	// handler if request body is invalid
	if err := c.BodyParser(&input); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	// update fields in the database
	update := bson.M{
		"$set": bson.M{
			"nama_kompetensi":  input.NamaKompetensi,
			"hard_skills":      input.HardSkills,
			"soft_skills":      input.SoftSkills,
			"model.updated_at": time.Now(),
		},
	}

	// update data in collection based on their "kompetensi_id" or params
	_, err = collection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return InternalServerError(c, "Failed to update competence data")
	}

	// return success
	return Ok(c, "Sucess edited Competence data", update)
}

func DeleteKompetensi(c *fiber.Ctx) error {
	// get kompetensi_id from params
	idParam := c.Params("id")

	// convert params to integer data type
	kompetensiID, err := strconv.Atoi(idParam)
	if err != nil {
		return BadRequest(c, "Invalid ID")
	}

	// connect to collection in MongoDB
	collection := config.MongoClient.Database("certificate-generator").Collection("competence")

	// make filter to find document based on kompetensi_id
	filter := bson.M{"kompetensi_id": kompetensiID}

	// find competence
	var competenceData bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&competenceData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Competence not found")
		}
		fmt.Println("MongoDB FindOne Error:", err)
		return InternalServerError(c, "Failed to fetch data")
	}

	// check if competence already deleted
	if deletedAt, ok := competenceData["model"].(bson.M)["deleted_at"]; ok && deletedAt != nil {
		// return the deletion time if the competence is already deleted
		return AlreadyDeleted(c, "This competence has already been deleted", deletedAt)
	}

	// make update for input timestamp DeletedAt
	update := bson.M{"$set": bson.M{"model.deleted_at": time.Now()}}

	// update document in collection MongoDB
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return InternalServerError(c, "Failed to delete competence")
	}

	// check if the document is found and updated
	if result.MatchedCount == 0 {
		return NotFound(c, "Competence not found")
	}

	// return success
	return Ok(c, "Sucess deleted Competence data", kompetensiID)
}

// function to get all kompetensi data
func GetAllKompetensi(c *fiber.Ctx) error {
	var results []bson.M

	collection := config.MongoClient.Database("certificate-generator").Collection("competence")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set the projection to return the required fields
	projection := bson.M{
		"_id":             1, // 0 to exclude the field
		"kompetensi_id":   1,
		"nama_kompetensi": 1, // 1 to include the field, _id will be included by default
	}

	// find the projection
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "No Competence found")
		}
		return InternalServerError(c, "Failed to fetch data")
	}
	defer cursor.Close(ctx)

	// decode each document and append it to results
	for cursor.Next(ctx) {
		var competence bson.M
		if err := cursor.Decode(&competence); err != nil {
			return InternalServerError(c, "Failed to decode data")
		}
		results = append(results, competence)
	}
	if err := cursor.Err(); err != nil {
		return InternalServerError(c, "Cursor error")
	}

	// return success
	return Ok(c, "Sucess get all Competence data", results)
}

// function to get detail kompetensi data based on their kompetensi_id
func GetDetailKompetensi(c *fiber.Ctx) error {
	// get kompetensi_id from params
	idParam := c.Params("id")

	// parsing kompetensi_id to integer type data
	kompetensiID, err := strconv.Atoi(idParam)
	if err != nil {
		return BadRequest(c, "Invalid ID")
	}

	// connect to collection in MongoDB
	collection := config.MongoClient.Database("certificate-generator").Collection("competence")

	// make filter to find document based on kompetensi_id (incremental id)
	filter := bson.M{"kompetensi_id": kompetensiID}

	// variable to hold search results
	var kompetensiDetail bson.M

	// find a single document that matches the filter
	err = collection.FindOne(context.TODO(), filter).Decode(&kompetensiDetail)
	if err != nil {
		// if not found, return a 404 status
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Data not found")
		}
		// if in server error, return status 500
		return InternalServerError(c, "Failed to retrieve data")
	}

	// check if document is already deleted
	if deletedAt, ok := kompetensiDetail["model"].(bson.M)["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the account is already deleted
		return AlreadyDeleted(c, "This competence has already been deleted", deletedAt)
	}

	// return success
	return Ok(c, "Sucess get detail Competence data", kompetensiDetail)
}
