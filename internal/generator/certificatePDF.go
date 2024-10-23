package generator

import (
	"html/template"
	"log"
	"os"
	model "pkl/finalProject/certificate-generator/model"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
)

func CreatePDF(c *fiber.Ctx, dataReq *model.CertificateData, zoom float64) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}

	if err := c.Render("temp/index.html", struct {
		Data model.CertificateData
		Enc  template.Srcset
	}{Data: *dataReq, Enc: template.Srcset(dataReq.QRCode.QRCodeEnc)}); err != nil { // remember to change index.html path
		return err
	}

	pdfTempl, err := os.Create("temp/temp" + dataReq.DataID + ".html")
	if err != nil {
		return err
	}
	defer func() {
		pdfTempl.Close()
		if err := os.Remove("temp/temp" + dataReq.DataID + ".html"); err != nil {
			log.Println("WARNING: memory leak (can't remove html file)\nerr:", err)
		}
	}()

	if _, err := pdfTempl.Write(c.Response().Body()); err != nil {
		return err
	}

	page := wkhtmltopdf.NewPage("temp/temp" + dataReq.DataID + ".html")
	page.Zoom.Set(zoom) // agak tergantung siapa yg buat file pdf-nya, di rendy zoom 1.064, di aku 1.3

	pdfg.AddPage(page)
	pdfg.MarginTop.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	err = pdfg.Create()
	if err != nil {
		return err
	}

	return pdfg.WriteFile("temp/certificate/" + dataReq.DataID + ".pdf")
}
