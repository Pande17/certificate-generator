package rest

import (
	"bytes"
	"certificate-generator/internal/generator"
	"certificate-generator/model"
	"fmt"
	"html/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
)

func CheckPDF(c *fiber.Ctx) error {
	var pdfReq struct {
		Data model.CertificateData `json:"data" bson:"data"`
		Type string                `json:"type"`
	}
	if err := c.BodyParser(&pdfReq); err != nil {
		return BadRequest(c, "Invalid body request", err.Error())
	}
	if !(pdfReq.Type == "a" || pdfReq.Type == "b") {
		return BadRequest(c, "Tipe sertifikat tidak diketahui.", "query type isn't a or b")
	}
	generator.CreatePDF(c, &pdfReq.Data, pdfReq.Type)
	return c.SendFile("assets/certificate/" + pdfReq.Data.DataID + "-" + pdfReq.Type + ".pdf")
}

func HandleBuildPdf(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles("assets/index.html")
	if err != nil {
		fmt.Printf("Failed to parse content file, err %s\n", err.Error())
		return BadRequest(c, "PDF Source Parse Failed, requested content not valid", err.Error())
	}

	var parser bytes.Buffer

	// If need to replace any content, do here in the map[string]any
	err = tmpl.Execute(&parser, map[string]any{})
	if err != nil {
		fmt.Printf("Failed to embed content to html file, err %s\n", err.Error())
		return BadRequest(c, "PDF Build Failed, requested content not valid", err.Error())
	}

	err = buildPDF(parser, "output/test.pdf")
	if err != nil {
		fmt.Printf("Failed to build pdf file, err %s\n", err.Error())
		return InternalServerError(c, "PDF Build Failed", err.Error())
	}

	return OK(c, "PDF Build Success", parser)
}

func buildPDF(parse bytes.Buffer, outputName string) error {
	pdf := wkhtmltopdf.NewPDFPreparer()
	res := wkhtmltopdf.NewPageReader(&parse)
	res.DisableExternalLinks.Set(true)
	res.EnableLocalFileAccess.Set(false)
	res.DisableInternalLinks.Set(false)

	// Set Header And Footer if Any
	// res.HeaderHTML.Set("assets/sample_header.html")
	// res.FooterHTML.Set("assets/sample_footer.html")

	pdf.AddPage(res)

	// Set PDF Margin
	pdf.MarginTop.Set(0)
	pdf.MarginLeft.Set(0)
	pdf.MarginBottom.Set(0)
	pdf.MarginRight.Set(0)
	res.Zoom.Set(1.367)

	js, err := pdf.ToJSON()

	if err != nil {
		return err
	}

	pdfFromJson, err1 := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewReader(js))
	if err1 != nil {
		return err
	}

	err = pdfFromJson.Create()
	if err != nil {
		return err
	}

	return pdfFromJson.WriteFile(outputName)
}
