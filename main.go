package main

import (
	"fmt"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// Fungsi untuk membuat PDF dari file HTML
func CreatePDF(htmlFilePath, outputPDFPath string, zoom float64) error {
	// Inisialisasi PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to create PDF generator: %w", err)
	}

	// Buat halaman baru dari file HTML yang akan dikonversi
	page := wkhtmltopdf.NewPage(htmlFilePath)
	page.Zoom.Set(zoom)

	// Set pengaturan PDF
	pdfg.AddPage(page)
	pdfg.MarginTop.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	page.EnableLocalFileAccess.Set(true) // Enable access to local resources
	page.DisableSmartShrinking.Set(true) // Optionally disable smart shrinking

	// Buat PDF dari halaman yang ditentukan
	err = pdfg.Create()
	if err != nil {
		return fmt.Errorf("failed to create PDF: %w", err)
	}

	// Tulis PDF ke file
	err = pdfg.WriteFile(outputPDFPath)
	if err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	fmt.Println("PDF berhasil dibuat:", outputPDFPath)
	return nil
}

func main() {
	// Path file HTML yang ingin dikonversi
	htmlFilePath := "temp/index.html"
	// Path untuk menyimpan output PDF
	outputPDFPath := "output/1.pdf"
	// Zoom level untuk konversi
	zoom := 1.21
	// Panggil fungsi CreatePDF untuk melakukan konversi
	err := CreatePDF(htmlFilePath, outputPDFPath, zoom)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}