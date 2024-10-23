package generator

import (
	"encoding/base64"
	"log"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(link, str string) (string, error) {
	qr, err := qrcode.New(link+str+".pdf", qrcode.Medium)
	if err != nil {
		return "", err
	}
	png, err := qr.PNG(-4)
	if err != nil {
		return "", err
	}
	encstr := base64.RawURLEncoding.EncodeToString(png)
	log.Println(encstr)
	return encstr, nil
}
