package generator

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(link string) (string, error) {
	qr, err := qrcode.New(link, qrcode.Medium)
	if err != nil {
		return "", err
	}
	png, err := qr.PNG(-4)
	if err != nil {
		return "", err
	}
	encstr := base64.StdEncoding.EncodeToString(png)
	return encstr, nil
}
