package rest

import (
	"log"
	"os"
	model "pkl/finalProject/certificate-generator/model"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
)

func TEMPlate(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON("nothing here yet")
}

func CreatePDF(c *fiber.Ctx) error {
	var dataReq struct {
		model.CertificateData `json:"data"`
		Zoom                  float64 `json:"zoom"`
	}
	if err := c.BodyParser(&dataReq); err != nil {
		return InternalServerError(c, err.Error(), "can't parse body req")
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return InternalServerError(c, "wkhtmltopdf not found", err.Error())
	}

	if err := c.Render("temp/index.html", dataReq.CertificateData); err != nil { // remember to change index.html path
		return InternalServerError(c, "can't render html", err.Error())
	}

	pdfTempl, err := os.Create("temp/temp" + dataReq.CertificateData.DataID + ".html")
	if err != nil {
		return InternalServerError(c, "can't create html file", err.Error())
	}
	defer func() {
		pdfTempl.Close()
		if err := os.Remove("temp/temp" + dataReq.CertificateData.DataID + ".html"); err != nil {
			log.Println("WARNING: memory leak (can't remove html file)\nerr:", err)
		}
	}()

	if _, err := pdfTempl.Write(c.Response().Body()); err != nil {
		return InternalServerError(c, "can't write to html file", err.Error())
	}

	page := wkhtmltopdf.NewPage("temp/temp" + dataReq.CertificateData.DataID + ".html")
	page.Zoom.Set(dataReq.Zoom) // agak tergantung siapa yg buat file pdf-nya, di rendy zoom 1.064, di aku 1.264

	pdfg.AddPage(page)
	pdfg.MarginTop.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	err = pdfg.Create()
	if err != nil {
		return InternalServerError(c, err.Error(), "can't create pdf")
	}

	err = pdfg.WriteFile("temp/" + dataReq.DataID + ".pdf")
	if err != nil {
		return InternalServerError(c, err.Error(), "can't writefile pdf")
	}
	return OK(c, "Success create certificate pdf", dataReq)
}
