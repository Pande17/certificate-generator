package rest

import (
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
)

func TEMPlate(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON("nothing here yet")
}

func CreatePDF(c *fiber.Ctx) error {
	var shshsj struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := c.BodyParser(&shshsj); err != nil {
		return InternalServerError(c, err.Error(), "can't parse body req")
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return InternalServerError(c, err.Error(), "wkhtmltopdf not found")
	}

	f, err := os.Open("temp/index.html")
	if err != nil {
		return InternalServerError(c, err.Error(), "can't find html file")
	}
	defer f.Close()

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))
	pdfg.MarginTop.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		return InternalServerError(c, err.Error(), "can't create pdf")
	}

	err = pdfg.WriteFile("temp/" + shshsj.ID + ".pdf")
	if err != nil {
		return InternalServerError(c, err.Error(), "can't writefile pdf")
	}
	return OK(c, "Success create certificate pdf", shshsj)
}
