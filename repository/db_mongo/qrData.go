package dbmongo

import (
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QrData struct {
	ID		primitive.ObjectID		`bson:"_id" json:"id"`
	PdfId	string					`bson:"pdfId" json:"pdf_id"`
	QrLink	string					`bson:"qrLink" json:"qr_link"`
	QrImage	base64.Encoding			`bson:"qrImage" json:"qr_image"`
	Model
}