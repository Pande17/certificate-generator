package controller

import (
	"context"
	"fmt"
	"pkl/finalProject/certificate-generator/repository/config"
	dbmongo "pkl/finalProject/certificate-generator/repository/db_mongo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateCompetence(c *fiber.Ctx) error {
	// struct for the incoming request body
	var competenceReq struct {
		ID				primitive.ObjectID 	`bson:"_id,omitempty"`
		CompetenceName	string				`bson:"competence_name"`
	}

	// parse the request body
	if err := c.BodyParser(&competenceReq); err != nil {
		return BadRequest(c, "Failed to read body")
	}

	// connect collection competence in database
	collection := config.MongoClient.Database("certificate-generator").Collection("competence")

	// new variable to check the availability of the competence name
	var existingCompetence dbmongo.Competence

	// new variable to find competence based on their name "competence_name"
	filter := bson.M{"competence_name": competenceReq.CompetenceName}

	// find competence with same competence name as input name
	err = collection.FindOne(context.TODO(), filter).Decode(&existingCompetence)
	if err == nil {
		return Conflict(c, "Competence name already exists")
	} else if err != mongo.ErrNoDocuments {
		return InternalServerError(c, "Error chechking for existing competence name")
	}

	nextCompetenceID, err := GetNextAccID()

	// return success
    return Ok(c, "Sucess created New Competence", nil)
}

func EditCompetence(c *fiber.Ctx) error {

	// return success
    return Ok(c, "Sucess edited Competence data", nil)
}

func DeleteCompetence(c *fiber.Ctx) error {

	// return success
    return Ok(c, "Sucess deleted Competence data", nil)
}

func SeeAllCompetence(c *fiber.Ctx) error {

	// return success
    return Ok(c, "Sucess get all Competence data", nil)
}

func SeeDetailCompetence(c *fiber.Ctx) error {

	// return success
    return Ok(c, "Sucess get Competence data", nil)
}

// Function for generate ID incremental for admin account
func GetNextCompetenceID(adminCollection *mongo.Collection) (int64, error) {
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