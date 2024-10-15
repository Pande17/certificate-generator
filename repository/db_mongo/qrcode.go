package dbmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QRCode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	QRCodeID    uint64             `bson:"qrcode_id"`
	QRCodePDFID string             `bson:"qrcode_str"`
	QRCodeLink  string             `bson:"qrcode_link"`
	Model
}
