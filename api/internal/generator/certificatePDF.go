package generator

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	model "pkl/finalProject/certificate-generator/model"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
)

var pdfg *wkhtmltopdf.PDFGenerator

// letter width is 4 for lowercase letters, and 8 for uppercase letters
//
// add values below to get that letter's width
//
// unk chars are considered 5
var txtWide = map[rune]int{
	'a': 0, 'A': 0,
	'b': 0, 'B': -1,
	'c': -1, 'C': -1,
	'd': 0, 'D': 1,
	'e': -1, 'E': -2,
	'f': -1, 'F': -2,
	'g': 0, 'G': -1,
	'h': 0, 'H': 0,
	'i': -2, 'I': -4,
	'j': -2, 'J': -3,
	'k': 0, 'K': 0,
	'l': -2, 'L': 1,
	'm': 1, 'M': 2,
	'n': 0, 'N': -1,
	'o': 0, 'O': -1,
	'p': 0, 'P': 0,
	'q': 0, 'Q': -1,
	'r': -1, 'R': 1,
	's': -1, 'S': -1,
	't': -1, 'T': -2,
	'u': 0, 'U': 0,
	'v': 0, 'V': -2,
	'w': 1, 'W': 1,
	'x': 0, 'X': -1,
	'y': 1, 'Y': 0,
	'z': 0, 'Z': 0,
	' ': -2,
}

type Funcs struct {
	Text string
}

func (s Funcs) SplitTxt(t string) []string {
	return strings.Split(t, "")
}

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

func CreatePDF(c *fiber.Ctx, dataReq *model.CertificateData, zoom float64, pageName string) error { //, pageNum string) error {
	type renderer struct {
		Data      model.CertificateData
		Enc       template.Srcset
		StyleReg  template.CSS
		StylePage template.CSS
		NameSize  string
		Funcs
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

	scale := float64(1)
	txtWidth := len(dataReq.NamaPeserta) * 4
	for _, c := range dataReq.NamaPeserta {
		if val, ok := txtWide[c]; ok {
			if string(c) == strings.ToUpper(string(c)) && c != ' ' {
				txtWidth += 4
			}
			txtWidth += val
		} else {
			txtWidth++
		}
	}
	if txtWidth > 120 {
		scale = 120 / float64(txtWidth)
	}
	scale = math.Floor(scale * 48)

	if err := c.Render("temp/index.html", renderer{
		Data:     *dataReq,
		Enc:      template.Srcset(dataReq.QRCode),
		StyleReg: template.CSS(stReg),
		//StylePage: template.CSS(stPage),
		NameSize: fmt.Sprintf("style='font-size:%.fpx'", scale),
		Funcs:    Funcs{Text: "foo"},
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
