package model

type QRCode struct {
	//QRCodeID    uint64 `json:"qrcode_id" bson:"qrcode_id"`
	QRCodePDFID string `json:"qrcode_str" bson:"qrcode_str"`
	QRCodeLink  string `json:"qrcode_link" bson:"qrcode_link"`
	QRCodeEnc   string `json:"qrcode_enc" bson:"qrcode_enc"`
}
