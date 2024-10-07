package controller

import (
	"context"
	"pkl/finalProject/certificate-generator/generator"
	"pkl/finalProject/certificate-generator/repository/config"
	dbmongo "pkl/finalProject/certificate-generator/repository/db_mongo"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCertificate(c *fiber.Ctx) error {
	// add body request
	var pdfReq struct {
		Data dbmongo.CertificateData `json:"data" bson:"data"`
	}

	// parse the body request
	if err := c.BodyParser(&pdfReq); err != nil {
		return BadRequest(c, "Invalid body request", "Invalid req body")
	}

	// connect collection certificate in database
	certificateCollection := config.GetCollection(config.MongoClient, "certificate")
	competenceCollection := config.GetCollection(config.MongoClient, "competence")
	counterCollection := config.GetCollection(config.MongoClient, "counters")

	// generate DataID (random string with 8 letter)
	newDataID, err := generator.GetUniqueRandomID(certificateCollection, 8)
	if err != nil {
		return InternalServerError(c, "Failed to generate Data ID", "Server failed generate Data ID")
	}

	// generate qrcode
	// newQrCode, err := GenerateQRCode("")

	now := time.Now()

	// generate referral ID (incremental ID)
	nextReferralID, err := generator.GenerateReferralID(counterCollection, now)
	if err != nil {
		return InternalServerError(c, "Failed to generate Referral ID", "Server failed generate Referral ID")
	}

	// fetch Kompetensi by the given nama_kompetensi from the request
	var kompetensi dbmongo.Kompetensi
	filter := bson.M{"nama_kompetensi": pdfReq.Data.Kompetensi}
	err = competenceCollection.FindOne(context.TODO(), filter).Decode(&kompetensi)
	if err != nil {
		return NotFound(c, "Competence Not Found", "Fetch Kompetepetensi by the given nama_kompetensi from the request")
	}

	hardSkills := kompetensi.HardSkills
	softSkills := kompetensi.SoftSkills

	mappedKompetensi := dbmongo.Kompetensi{
		ID: primitive.NewObjectID(),
		KompetensiID: kompetensi.KompetensiID,
		HardSkills: hardSkills,
		SoftSkills: softSkills,
	}


	// // map Kompetensi's hard & soft skills to CertificateData
	// var hardSkills []dbmongo.HardSkill
	// for _, hardSkill := range pdfReq.Data.HardSkillPDF {
	// 	descriptions := []dbmongo.Description{
	// 		{
	// 			UnitCode:  hardSkill.HardSkillCode,
	// 			UnitTitle: hardSkill.HardSkillDesc,
	// 		},
	// 	}
	// 	hardSkills = append(hardSkills, dbmongo.HardSkill{
	// 		HardSkillName: hardSkill.HardSkillName,
	// 		Descriptions:  descriptions,
	// 	})
	// }

	// var softSkills []dbmongo.SoftSkill
	// for _, softSkill := range pdfReq.Data.SoftSkillPDF {
	// 	descriptions := []dbmongo.Description{
	// 		{
	// 			UnitCode:  softSkill.SoftSkillCode,
	// 			UnitTitle: softSkill.SoftSkillDesc,
	// 		},
	// 	}
	// 	softSkills = append(softSkills, dbmongo.SoftSkill{
	// 		SoftSkillName: softSkill.SoftSkillName,
	// 		Descriptions:  descriptions,
	// 	})
	// }

	// // create the Kompetensi Object
	// kompetensi := dbmongo.Kompetensi{
	// 	ID:             primitive.NewObjectID(),
	// 	KompetensiID:   uint64(nextReferralID),
	// 	NamaKompetensi: pdfReq.Data.Kompetensi,
	// 	HardSkills:     hardSkills,
	// 	SoftSkills:     softSkills,
	// }

	certificate := dbmongo.PDF{
		ID:     primitive.NewObjectID(),
		DataID: newDataID,
		Data: dbmongo.CertificateData{
			SertifName: pdfReq.Data.SertifName,
			KodeReferral: dbmongo.KodeReferral{
				ReferralID: nextReferralID,
				Divisi:     pdfReq.Data.KodeReferral.Divisi,
				BulanRilis: pdfReq.Data.KodeReferral.BulanRilis,
				TahunRilis: pdfReq.Data.KodeReferral.TahunRilis,
			},
			NamaPeserta:    pdfReq.Data.NamaPeserta,
			SKKNI:          pdfReq.Data.SKKNI,
			KompetenBidang: pdfReq.Data.KompetenBidang,
			Kompetensi:     mappedKompetensi.NamaKompetensi,
			Validation:     pdfReq.Data.Validation,
			// KodeQR:         newQrCode,
			DataID: newDataID,
			// TotalJP:   pdfReq.Data.TotalJP,
			TotalMeet:    pdfReq.Data.TotalMeet,
			MeetTime:     pdfReq.Data.MeetTime,
			ValidDate:    pdfReq.Data.ValidDate,
			HardSkillPDF: mappedKompetensi.HardSkills,
			SoftSkillPDF: mappedKompetensi.SoftSkills,
			// FinalSkor:      float64(pdfReq.Data.FinalSkor),
		},
		Model: dbmongo.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// insert data from struct "PDF" to collection "certificate" in database MongoDB
	_, err = certificateCollection.InsertOne(context.TODO(), certificate)
	if err != nil {
		return InternalServerError(c, "Failed to create new certificate data", "Server failed create new certificate")
	}

	// return success
	return OK(c, "Success create new certificate", certificate)
}

// {
//     "sertif_name": "Sertifikat pertama",
//     "kode_referral": [
//         {
//             "divisi": "BIS",
//             "bulan_rilis": "V",
//             "tahun_rilis": "2024"
//         }
//     ],
//     "nama_peserta": "I Kadek Pande Feri Dwi Wijaya",
//     "skkni": "SKKNI Nomor 56 Tahun 2018",
//     "kompeten_bidang": "Pengembangan Bisnis",
//     "kompetensi": "Leadership & Building Startup",
//     "valid_date": [
//         {
//             "valid_total": "3 Tahun",
//             "valid_start": "29 Agustus 2024",
//             "valid_end": "29 Agustus 2027"
//         }
//     ],
//     "validation": "Denpasar, 29 Agustus 2024",
//     "total_meet": "14 Pertemuan",
//     "meet_time": "3.5 Bulan"
// }
