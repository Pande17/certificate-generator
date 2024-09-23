package controller

import "github.com/skip2/go-qrcode"

func GenerateQRCode(str string) ([]byte, error) {
	return qrcode.Encode(str, qrcode.Medium, 256)
}
