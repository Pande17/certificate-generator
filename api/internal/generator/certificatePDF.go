package generator

import (
	"certificate-generator/model"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
)

var pdfg *wkhtmltopdf.PDFGenerator
var mtx sync.Mutex
var CreatingPDF map[string]bool

// functions for rendering html certificate
var funcs = template.FuncMap{
	"splittxt": func(s string) []string {
		return strings.Split(s, "")
	},
	"add": func(a, b int) int {
		return a + b
	},
	"rangecheck": func(s string) int {
		scale := float64(1)
		txtWidth := len(s) * 4
		for _, c := range s {
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
		return int(math.Floor(scale * 48))
	},
	"parity": func(i int) int {
		return i % 2
	},
}

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

func init() {
	var err error
	pdfg, err = wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println("wkhtmltopdf not found")
	}

	pdfg.MarginTop.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
}

func CreatePDF(c *fiber.Ctx, dataReq *model.CertificateData, certifType string) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Minute)
	defer cancel()
	return createPDF(ctx, c, dataReq, certifType)
}

func createPDF(ctx context.Context, c *fiber.Ctx, dataReq *model.CertificateData, certifType string) error {
	makeA := strings.Contains(certifType, "a")
	makeB := strings.Contains(certifType, "b")
	if !(makeA || makeB) {
		return errors.New("CreatePDF: certifType isn't a or b")
	}

	// generate qrcode
	link := fmt.Sprintf("%s://%s/assets/certificate/", c.Protocol(), c.Hostname())
	encstr, err := GenerateQRCode(link, dataReq.DataID)
	if err != nil {
		return err
	}
	dataReq.QRCode = encstr

	mtx.Lock()
	CreatingPDF[dataReq.DataID+"-a"] = makeA
	CreatingPDF[dataReq.DataID+"-b"] = makeB
	defer func() {
		CreatingPDF = map[string]bool{}
	}()
	defer mtx.Unlock()

	select {
	case <-ctx.Done():
		return errors.New("CreatePDF: timeout exceeded")
	default:
		break
	}

	page1, err := makePage(c, dataReq, "page1")
	if err != nil {
		return err
	} else {
		defer removeFile(page1)
	}
	page2a, err := func() (*wkhtmltopdf.Page, error) {
		if !makeA {
			return nil, nil
		}
		return makePage(c, dataReq, "page2a")
	}()
	if err != nil {
		return err
	} else if page2a != nil {
		defer removeFile(page2a)
	}
	page2b, err := func() (*wkhtmltopdf.Page, error) {
		if !makeB {
			return nil, nil
		}
		return makePage(c, dataReq, "page2b")
	}()
	if err != nil {
		return err
	} else if page2b != nil {
		defer removeFile(page2b)
	}

	if makeA {
		pdfg.ResetPages()
		pdfg.AddPage(page1)
		pdfg.AddPage(page2a)
		if err := pdfg.Create(); err != nil {
			return err
		}
		if err := pdfg.WriteFile("assets/certificate/" + dataReq.DataID + "-a.pdf"); err != nil {
			return err
		}
	}

	if makeB {
		pdfg.ResetPages()
		pdfg.AddPage(page1)
		pdfg.AddPage(page2b)
		if err := pdfg.Create(); err != nil {
			return err
		}
		if err := pdfg.WriteFile("assets/certificate/" + dataReq.DataID + "-b.pdf"); err != nil {
			return err
		}
	}

	return nil
}

func makePage(c *fiber.Ctx, dataReq *model.CertificateData, pageName string) (*wkhtmltopdf.Page, error) {
	t, err := template.New("").Funcs(funcs).ParseFiles("assets/"+pageName+".html", "assets/style.html")
	if err != nil {
		return nil, err
	}
	if err := t.ExecuteTemplate(c.Response().BodyWriter(), pageName, *dataReq); err != nil {
		return nil, err
	}

	fileName := "temp/temp" + dataReq.DataID + pageName + ".html"
	htmlSertif, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		htmlSertif.Close()
	}()

	if _, err := htmlSertif.Write(c.Response().Body()); err != nil {
		return nil, err
	}
	c.Response().Reset()

	page := wkhtmltopdf.NewPage(fileName)
	return page, nil
}

func removeFile(page *wkhtmltopdf.Page) {
	if err := os.Remove(page.InputFile()); err != nil {
		log.Println("WARNING: memory leak (can't remove html file)\nerr:", err)
	}
}
