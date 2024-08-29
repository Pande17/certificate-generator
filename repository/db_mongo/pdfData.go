package dbmongo

import (
	"encoding/base64"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File struct untuk tabel file
type PdfData struct {
	ID              primitive.ObjectID 	`bson:"_id,omitempty" json:"id"`
	UserId          primitive.ObjectID 	`bson:"userId" json:"user_id"`
	PdfId			string			   	`bson:"pdfId" json:"pdf_id"`
	QrId            string             	`bson:"qrId" json:"qr_id"`
	SertifName      string             	`bson:"sertifName" json:"sertif_name"`
	SertifId        uint               	`bson:"sertifId" json:"sertif_id"`
	BulanRilis      string             	`bson:"bulanRilis" json:"bulan_rilis"`
	TahunRilis      uint               	`bson:"tahunRilis" json:"tahun_rilis"`
	NamaPeserta     string             	`bson:"namaPeserta" json:"nama_peserta"`
	ProgramLat	    string			   	`bson:"programLat" json:"program_lat"`
	BidangKompeten  string             	`bson:"bidangKompeten" json:"bidang_kompeten"`
	DataNilai       primitive.ObjectID 	`bson:"dataNilai" json:"data_nilai"`
	ValidTotal      string             	`bson:"validTotal" json:"valid_total"`
	ValidDateStart  time.Time          	`bson:"validDateStart" json:"valid_date_start"`
	ValidDateEnd    time.Time          	`bson:"validDateEnd" json:"valid_date_end"`
	NilaiAkhir      float64            	`bson:"nilaiAkhir" json:"nilai_akhir"`
	PdfFile			base64.Encoding		`bson:"pdfFile" json:"pdf_file"`
	Model
}