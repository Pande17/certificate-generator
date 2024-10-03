package controller

import (
	"context"
	"pkl/finalProject/certificate-generator/repository/config"
	dbmongo "pkl/finalProject/certificate-generator/repository/db_mongo"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// LINK STRING NEEDS FIXING (still example only)
func GenerateQRCode(link, str string) error {
	return qrcode.WriteFile(link+str, qrcode.Medium, 256, "temp/"+str+".png")
}

// func for create new kompetensi
func CreateQRCode(c *fiber.Ctx) error {
	// struct for the incoming request body
	var qrCodeReq struct {
		ID          primitive.ObjectID `bson:"_id,omitempty"`
		QRCodeID    uint64             `bson:"qrcode_id"`
		QRCodePDFID string             `json:"pdf_id"`
	}

	// parse request body
	if err := c.BodyParser(&qrCodeReq); err != nil {
		return BadRequest(c, "Failed to read body")
	}

	link := "https://example.com/"

	if err := GenerateQRCode(link, qrCodeReq.QRCodePDFID); err != nil {
		return InternalServerError(c, "can't create qr code")
	}

	// get collection in db
	collection := config.MongoClient.Database("certificate-generator").Collection("qrCode")

	// check existing qrcode w/ the same str
	var existingQRCode dbmongo.QRCode
	filter := bson.M{"string_qrcode": qrCodeReq.QRCodePDFID}
	err := collection.FindOne(context.TODO(), filter).Decode(&existingQRCode)
	if err == nil {
		return Conflict(c, "QR Code already exists")
	} else if err != mongo.ErrNoDocuments {
		return InternalServerError(c, "Error checking for existing QR Code")
	}

	// generate incremental qrcode_id
	nextQRCodeID, err := GetNextIncrementalID(collection, "qrcode_id")
	if err != nil {
		return InternalServerError(c, "Failed to generate QR Code ID")
	}

	// struct to input data to db
	qrCodein := dbmongo.QRCode{
		ID:          primitive.NewObjectID(),
		QRCodeID:    uint64(nextQRCodeID),
		QRCodePDFID: qrCodeReq.QRCodePDFID,
		QRCodeLink:  link + qrCodeReq.QRCodePDFID,
		Model: dbmongo.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// insert data to db
	_, err = collection.InsertOne(context.TODO(), qrCodein)
	if err != nil {
		return InternalServerError(c, "Failed to create new QR Code")
	}

	// return success
	return Ok(c, "Success creating new QR Code", qrCodein)
}
