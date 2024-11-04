package generator

import (
	"html/template"
	"log"
	"os"
	model "pkl/finalProject/certificate-generator/model"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
)

var pdfg *wkhtmltopdf.PDFGenerator

func init() {
	var err error
	pdfg, err = wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println("wkhtmltopdf not found") // change to log.Fatal later
	}

	pdfg.MarginTop.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
}

func CreatePDF(c *fiber.Ctx, dataReq *model.CertificateData, zoom float64) error { //, pageNum string) error {
	type renderer struct {
		Data      model.CertificateData
		Enc       template.Srcset
		StyleReg  template.CSS
		StylePage template.CSS
	}

	pdfg.ResetPages()

	stReg, err := readCSS("styleReg")
	if err != nil {
		return err
	}
	// stPage, err := readCSS("stylePage" + pageNum)
	// if err != nil {
	// 	return err
	// }

	if err := c.Render("temp/index.html", renderer{
		Data:     *dataReq,
		Enc:      template.Srcset(dataReq.QRCode),
		StyleReg: template.CSS(stReg),
		//StylePage: template.CSS(stPage),
	}); err != nil {
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
	page.Zoom.Set(zoom)
	pdfg.AddPage(page)

	err = pdfg.Create()
	if err != nil {
		return err
	}

	return pdfg.WriteFile("temp/certificate/" + dataReq.DataID + ".pdf")
}

func readCSS(filename string) (string, error) {
	file, err := os.Open("temp/" + filename + ".html")
	if err != nil {
		return "", err
	}
	defer file.Close()

	fStat, err := file.Stat()
	if err != nil {
		return "", err
	}

	data := make([]byte, fStat.Size()*2)
	if _, err := file.Read(data); err != nil {
		return "", err
	}
	return string(data), nil
}
