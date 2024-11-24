package rest

import (
	"certificate-generator/database"
	"certificate-generator/internal/generator"
	"certificate-generator/model"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateCertificate(c *fiber.Ctx) error {
	// add body request
	var pdfReq struct {
		Data   model.CertificateData `json:"data" bson:"data"`
		SaveDB bool                  `json:"savedb"`
	}

	// parse the body request
	if err := c.BodyParser(&pdfReq); err != nil {
		return BadRequest(c, "Invalid body request", err.Error())
	}

	// connect collection certificate in database
	certificateCollection := database.GetCollection("certificate")
	competenceCollection := database.GetCollection("competence")
	counterCollection := database.GetCollection("counters")

	// generate DataID (random string with 8 letter)
	newDataID, err := generator.GetUniqueRandomID(certificateCollection, 8)
	if err != nil {
		return InternalServerError(c, "Failed to generate Data ID", "Server failed generate Data ID")
	}

	// generate referral ID
	currentTime := time.Now()
	nextReferralID, err := generator.GenerateReferralID(counterCollection, currentTime)
	if err != nil {
		return InternalServerError(c, "Failed to generate Referral ID", "Server failed generate Referral ID")
	}

	// generate month roman and year
	year := currentTime.Year()
	monthRoman := generator.MonthToRoman(int(currentTime.Month()))

	// fetch Kompetensi by the given nama_kompetensi from the request
	var kompetensi model.Kompetensi
	filter := bson.M{"nama_kompetensi": pdfReq.Data.Kompetensi}
	err = competenceCollection.FindOne(context.TODO(), filter).Decode(&kompetensi)
	if err != nil {
		return NotFound(c, "Competence Not Found", "Fetch Kompetepetensi by the given nama_kompetensi from the request")
	}

	totalHSJP, totalHSSkor := uint64(0), float64(0)
	for _, hs := range pdfReq.Data.HardSkills.Skills {
		totalHSJP += hs.SkillJP
		totalHSSkor += hs.SkillScore
	}

	totalSSJP, totalSSSkor := uint64(0), float64(0)
	for _, ss := range pdfReq.Data.SoftSkills.Skills {
		totalSSJP += ss.SkillJP
		totalSSSkor += ss.SkillScore
	}

	mappedData := model.CertificateData{
		SertifName: strings.ToUpper(pdfReq.Data.SertifName),
		KodeReferral: model.KodeReferral{
			ReferralID: nextReferralID,
			Divisi:     kompetensi.Divisi,
			BulanRilis: monthRoman,
			TahunRilis: year,
		},
		NamaPeserta: strings.TrimSpace(pdfReq.Data.NamaPeserta),
		SKKNI:       pdfReq.Data.SKKNI,
		Kompetensi:  pdfReq.Data.Kompetensi,
		Validation:  pdfReq.Data.Validation,
		DataID:      newDataID,
		TotalJP:     totalHSJP + totalSSJP,
		TotalMeet:   pdfReq.Data.TotalMeet,
		MeetTime:    pdfReq.Data.MeetTime,
		ValidDate:   pdfReq.Data.ValidDate,
		HardSkills: model.SkillPDF{
			Skills:          pdfReq.Data.HardSkills.Skills,
			TotalSkillJP:    totalHSJP,
			TotalSkillScore: float64(math.Round(totalHSSkor/float64(len(pdfReq.Data.HardSkills.Skills))*10) / 10),
		},
		SoftSkills: model.SkillPDF{
			Skills:          pdfReq.Data.SoftSkills.Skills,
			TotalSkillJP:    totalSSJP,
			TotalSkillScore: float64(math.Round(totalSSSkor/float64(len(pdfReq.Data.SoftSkills.Skills))*10) / 10),
		},
		FinalSkor: float64(math.Round((totalHSSkor+totalSSSkor)/float64(len(pdfReq.Data.HardSkills.Skills)+len(pdfReq.Data.SoftSkills.Skills))*10) / 10),
	}

	certificate := model.PDF{
		DataID:     newDataID,
		SertifName: strings.ToUpper(pdfReq.Data.SertifName),
		Data:       mappedData,
		Model: model.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
	}

	// make pdf creation concurrent to return handler faster
	go generator.CreatePDF(c, &mappedData, "ab")

	// insert data from struct "PDF" to collection "certificate" in database MongoDB
	if pdfReq.SaveDB {
		_, err = certificateCollection.InsertOne(context.TODO(), certificate)
		if err != nil {
			return InternalServerError(c, "Failed to create new certificate data", "Server failed create new certificate")
		}
	}

	// return success
	return OK(c, "Success create new certificate", certificate)
}

// function to get all kompetensi data
func GetAllCertificates(c *fiber.Ctx) error {
	var results []bson.M

	collection := database.GetCollection("certificate")
	ctx := c.Context()

	// set the projection to return the required fields
	projection := bson.M{
		"_id":         1,
		"data_id":     1,
		"sertif_name": 1,
		"created_at":  1,
		"updated_at":  1,
		"deleted_at":  1,
	}

	// find the projection
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "No certificate found", err.Error())
		}
		return InternalServerError(c, "Failed to fetch data", err.Error())
	}
	defer cursor.Close(ctx)

	// decode each document and append it to results
	for cursor.Next(ctx) {
		var certiticate bson.M
		if err := cursor.Decode(&certiticate); err != nil {
			return InternalServerError(c, "Failed to decode data", err.Error())
		}
		if deletedAt, ok := certiticate["deleted_at"]; ok && deletedAt != nil {
			// skip deleted certificates
			continue
		}
		results = append(results, certiticate)
	}
	if err := cursor.Err(); err != nil {
		return InternalServerError(c, "Cursor error", err.Error())
	}

	// return success
	return OK(c, "Sucess get all Certificate data", results)
}

func GetCertificateByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	searchKey := c.Params("type")
	if searchKey == "" { // from handler w/o type param, to not break api
		searchKey = "oid"
	} else if searchKey == "a" || searchKey == "b" { // from type of certificate
		searchKey = "data_id"
	}
	var searchVal any
	searchVal = idParam

	// Convert to ObjectID if needed
	if searchKey == "oid" {
		searchKey = "_id"
		certifID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			return BadRequest(c, "Sertifikat ini tidak ada!", "Please provide a valid ObjectID")
		}
		searchVal = certifID
	}

	// connect to collection in mongoDB
	collection := database.GetCollection("certificate")

	// make filter to find document based on data_id (incremental id)
	filter := bson.M{searchKey: searchVal}

	// variable to hold search results
	var certifDetail bson.M

	// find a single document that matches the filter
	if err := collection.FindOne(c.Context(), filter).Decode(&certifDetail); err != nil {
		// if not found, return a 404 status
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Data not found", "Find Detail Certificate")
		}
		// if in server error, return status 500
		return InternalServerError(c, "Failed to retrieve data", err.Error())
	}

	// Check if DeletedAt field already has a value
	if deletedAt, ok := certifDetail["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the certificate is already deleted
		return AlreadyDeleted(c, "This certificate has already been deleted", "Check deleted certificate", deletedAt)
	}

	// return success
	return OK(c, "Berhasil mendapatkan detail sertifikat.", certifDetail)
}

// Function for soft delete admin account
func DeleteCertificate(c *fiber.Ctx) error {
	// Get dataid from params
	idParam := c.Params("id")

	// Convert idParam to ObjectID if needed
	certifID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Sertifikat ini tidak ada!", "Please provide a valid ObjectID")
	}

	// connect to collection in mongoDB
	certificateCollection := database.GetCollection("certificate")

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"_id": certifID}

	// find admin account
	var certificate bson.M
	err = certificateCollection.FindOne(context.TODO(), filter).Decode(&certificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("error: %v\n", err.Error())
			return NotFound(c, "Certificate not found", "Cannot find certificate")
		}
		return InternalServerError(c, "Failed to fetch certificate", "server error cannot find certificate")
	}

	// Check if DeletedAt field already has a value
	if deletedAt, ok := certificate["deleted_at"]; ok && deletedAt != nil {
		// Return the deletion time if the certificate is already deleted
		return AlreadyDeleted(c, "This certificate has already been deleted", "Check deleted certificate", deletedAt)
	}

	// make update for input timestamp DeletedAt
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	// update document in collection MongoDB
	result, err := certificateCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return InternalServerError(c, "Failed to delete certificate", "Delete certificate")
	}

	// Check if the document is found and updated
	if result.MatchedCount == 0 {
		return NotFound(c, "Certificate not found", "Found certificate")
	}

	// Respons success
	return OK(c, "Successfully deleted certificate", idParam)
}

func DownloadCertificate(c *fiber.Ctx) error {
	certifType := c.Params("type")
	if !(certifType == "a" || certifType == "b") {
		return BadRequest(c, "Tipe sertifikat tidak diketahui.", "query type isn't a or b")
	}

	// search certif data, checking if data exists
	if c.Next(); c.Response().StatusCode()/100 != 2 {
		return NotFound(c, "Sertifikat dengan id "+c.Params("id", "yang dicari")+" tidak ditemukan.", "use certif that exists in db")
	}

	var resp fiber.Map
	var pdf model.PDF
	bodyres := c.Response().Body()
	if err := json.Unmarshal(bodyres, &resp); err != nil {
		return InternalServerError(c, "eror", err.Error())
	}
	if pdfBytes, err := json.Marshal(resp["data"]); err != nil {
		return InternalServerError(c, "eror", err.Error())
	} else {
		if err := json.Unmarshal(pdfBytes, &pdf); err != nil {
			return InternalServerError(c, "eror", err.Error())
		}
	}
	data := pdf.Data

	filepath := "./assets/certificate/" + data.DataID + "-" + certifType + ".pdf"
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			if _, creating := generator.CreatingPDF[data.DataID+"-"+certifType]; creating {
				for _, creating := generator.CreatingPDF[data.DataID+"-"+certifType]; creating; {
					time.Sleep(time.Second)
				}
			} else {
				if err = generator.CreatePDF(c, &data, certifType); err != nil {
					return InternalServerError(c, "can't create pdf file", err.Error())
				}
			}
		} else {
			return InternalServerError(c, "eror", err.Error())
		}
	}
	c.Response().Header.Add("Content-Type", "application/pdf")
	return c.Download("./assets/certificate/"+data.DataID+"-"+certifType+".pdf", "Sertifikat BTW Edutech "+certifType+" - "+data.NamaPeserta)
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
