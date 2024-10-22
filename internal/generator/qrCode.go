package generator

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(link, str string) (string, error) {
	encstr := base64.URLEncoding.EncodeToString([]byte(link + str))
	return encstr, qrcode.WriteFile(link+str, qrcode.Medium, 256, "temp/"+str+".png")
}
