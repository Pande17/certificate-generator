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
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect collection certificate in database
var certificateCollection = database.GetCollection("certificate")
var competenceCollection = database.GetCollection("competence")
var counterCollection = database.GetCollection("counters")

// function for Create data certificate
func CreateCertificate(c *fiber.Ctx) error {
	// add body request
	var pdfReq struct {
		Data model.CertificateData `json:"data" bson:"data"`
	}

	// parse the body request
	if err := c.BodyParser(&pdfReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Mohon periksa kembali.", err.Error())
	}

	// Retrieve the admin ID from the claims stored in context
	claims := c.Locals("admin").(jwt.MapClaims)
	adminID, ok := claims["sub"].(string)
	if !ok {
		return Unauthorized(c, "Token Admin tidak valid!", "Token Admin tidak valid!")
	}

	// Convert adminID (which is a string) to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return Unauthorized(c, "Format token admin tidak valid!", "Format token admin tidak valid!")
	}

	// generate DataID (random string with 8 letter)
	newDataID, err := generator.GetUniqueRandomID(certificateCollection, 8)
	if err != nil {
		return Conflict(c, "Gagal membuat ID Sertifikat! Silahkan coba lagi.", "Server failed generate Data ID")
	}

	// generate referral ID
	currentTime := time.Now()
	nextReferralID, err := generator.GenerateReferralID(counterCollection, currentTime)
	if err != nil {
		return Conflict(c, "Gagal membuat sertifikat! Silahkan coba lagi.", "Server failed generate Referral ID")
	}

	// generate month roman and year
	year := currentTime.Year()
	monthRoman := generator.MonthToRoman(int(currentTime.Month()))

	// fetch Kompetensi by the given nama_kompetensi from the request
	var kompetensi model.Kompetensi
	filter := bson.M{"nama_kompetensi": pdfReq.Data.Kompetensi}
	err = competenceCollection.FindOne(context.TODO(), filter).Decode(&kompetensi)
	if err != nil {
		return NotFound(c, "Gagal memeriksa kompetensi yang ada. Silakan coba lagi.", err.Error())
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

	sertifName := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.ToUpper(pdfReq.Data.SertifName)), "SERTIFIKAT"))
	mappedData := model.CertificateData{
		AdminId:    objectID,
		SertifName: sertifName,
		Logo:       pdfReq.Data.Logo,
		KodeReferral: model.KodeReferral{
			ReferralID: nextReferralID,
			Divisi:     kompetensi.Divisi,
			BulanRilis: monthRoman,
			TahunRilis: year,
		},
		NamaPeserta:    pdfReq.Data.NamaPeserta,
		SKKNI:          pdfReq.Data.SKKNI,
		KompetenBidang: pdfReq.Data.KompetenBidang,
		Kompetensi:     pdfReq.Data.Kompetensi,
		Validation:     pdfReq.Data.Validation,
		DataID:         newDataID,
		TotalJP:        totalHSJP + totalSSJP,
		TotalMeet:      pdfReq.Data.TotalMeet,
		MeetTime:       pdfReq.Data.MeetTime,
		ValidDate:      pdfReq.Data.ValidDate,
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
		Signature: model.Signature{
			ConfigName: pdfReq.Data.Signature.ConfigName,
			Stamp:      pdfReq.Data.Signature.Stamp,
			Signature:  pdfReq.Data.Signature.Signature,
			Name:       pdfReq.Data.Signature.Name,
			Role:       pdfReq.Data.Signature.Role,
		},
	}

	certificate := model.PDF{
		AdminId:    objectID,
		DataID:     newDataID,
		SertifName: sertifName,
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
	_, err = certificateCollection.InsertOne(context.TODO(), certificate)
	if err != nil {
		return Conflict(c, "Gagal membuat data sertifikat baru! Silakan coba lagi.", "Server failed create new certificate")
	}

	// return success
	return OK(c, "Berhasil membuat sertifikat baru!", certificate)
}

// function for get all certificate data
func GetAllCertificates(c *fiber.Ctx) error {
	var results []bson.M

	ctx := c.Context()

	// Retrieve the admin ID from the claims stored in context
	claims := c.Locals("admin").(jwt.MapClaims)
	adminID, ok := claims["sub"].(string)
	if !ok {
		return Unauthorized(c, "Token Admin tidak valid!", "Token Admin tidak valid!")
	}

	// Convert adminID (which is a string) to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return Unauthorized(c, "Format token admin tidak valid!", "Format token admin tidak valid!")
	}

	// set the projection to return the required fields
	projection := bson.M{
		"_id":         1,
		"admin_id":    1,
		"data_id":     1,
		"sertif_name": 1,
		"created_at":  1,
		"updated_at":  1,
		"deleted_at":  1,
	}

	// Create the filter to include admin_id and handle deleted_at
	filter := bson.M{
		"admin_id": objectID,
		"$or": []bson.M{
			{"deleted_at": bson.M{"$exists": false}}, // DeletedAt field does not exist
			{"deleted_at": bson.M{"$eq": nil}},       // DeletedAt field is nil
		},
	}

	// Find all documents that match
	cursor, err := certificateCollection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Sertifikat tidak dapat ditemukan!", err.Error())
		}
		return Conflict(c, "Gagal mengambil data kompetensi! Silakan coba lagi.", err.Error())
	}
	defer cursor.Close(ctx)

	// Decode each document and append it to results
	for cursor.Next(ctx) {
		var certiticate bson.M
		if err := cursor.Decode(&certiticate); err != nil {
			return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", err.Error())
		}
		results = append(results, certiticate)
	}
	if err := cursor.Err(); err != nil {
		return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", err.Error())
	}

	// return success
	return OK(c, "Berhasil menampilkan semua data Kompetensi!", results)
}

// function for get detail certificate data by id
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
			return BadRequest(c, "Sertifikat ini tidak ada!", "Please provide a valid ObjectID")
		}
		searchVal = certifID
	}

	// Make filter to find document based on search key & value
	filter := bson.M{searchKey: searchVal}

	// Variable to hold search results
	var certifDetail bson.M

	// find a single document that matches the filter
	if err := certificateCollection.FindOne(c.Context(), filter).Decode(&certifDetail); err != nil {
		// if not found, return a 404 status
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Sertifikat ini tidak dapat ditemukan! Silakan periksa ID yang dimasukkan.", "Find Detail Certificate")
		}
		return Conflict(c, "Gagal mendapatkan data! Silakan coba lagi.", err.Error())
	}

	// Check if the certificate has been deleted
	if deletedAt, exists := certifDetail["deleted_at"]; exists && deletedAt != nil {
		return AlreadyDeleted(c, "Sertifikat ini telah dihapus! Silakan hubungi admin.", "Check deleted certificate", deletedAt)
	}

	// return success
	return OK(c, "Berhasil mendapatkan data sertifikat!", certifDetail)
}

// / Function for edit data certificate
func EditCertificate(c *fiber.Ctx) error {
	idParam := c.Params("id")
	accID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Data Sertifikat tidak ditemukan! Silakan periksa ID sertifikat.", "Gagal mengonversi parameter pada Edit Sertifikat")
	}

	filter := bson.M{"_id": accID}
	var certificate bson.M

	if err := certificateCollection.FindOne(c.Context(), filter).Decode(&certificate); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Data Sertifikat tidak ditemukan!", "Gagal menemukan Data Sertifikat")
		}
		return Conflict(c, "Gagal mendapatkan Data Sertifikat! Silakan coba lagi.", "Gagal menemukan Data Sertifikat")
	}

	if deletedAt, ok := certificate["deleted_at"]; ok && deletedAt != nil {
		return AlreadyDeleted(c, "Data Sertifikat ini sudah dihapus! Silakan hubungi Data Sertifikat.", "Periksa Data Sertifikat yang dihapus", deletedAt)
	}

	var input struct {
		Data struct {
			SertifName   string `json:"sertif_name"`
			Logo         string `json:"logo"`
			KodeReferral struct {
				Divisi string `json:"divisi"`
			} `json:"kode_referral"`
			NamaPeserta    string `json:"nama_peserta"`
			Skkni          string `json:"skkni"`
			KompetenBidang string `json:"kompeten_bidang"`
			Kompetensi     string `json:"kompetensi"`
			Validation     string `json:"validation"`
			ValidDate      struct {
				ValidTotal string `json:"valid_total"`
				ValidStart string `json:"valid_start"`
				ValidEnd   string `json:"valid_end"`
			} `json:"valid_date"`
			TotalMeet  int    `json:"total_meet"`
			MeetTime   string `json:"meet_time"`
			HardSkills struct {
				Skills []struct {
					SkillName   string `json:"skill_name"`
					Description []struct {
						UnitCode  string `json:"unit_code"`
						UnitTitle string `json:"unit_title"`
					} `json:"description"`
					SkillJP    int     `json:"skill_jp"`
					SkillScore float64 `json:"skill_score"`
				} `json:"skills"`
			} `json:"hard_skills"`
			SoftSkills struct {
				Skills []struct {
					SkillName   string `json:"skill_name"`
					Description []struct {
						UnitCode  string `json:"unit_code"`
						UnitTitle string `json:"unit_title"`
					} `json:"description"`
					SkillJP    int     `json:"skill_jp"`
					SkillScore float64 `json:"skill_score"`
				} `json:"skills"`
			} `json:"soft_skills"`
			Signature struct {
				ConfigName string `json:"config_name"`
				Stamp      string `json:"stamp"`
				Signature  string `json:"signature"`
				Name       string `json:"name"`
				Role       string `json:"role"`
			} `json:"signature"`
		} `json:"data"`
	}

	if err := c.BodyParser(&input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Silakan periksa kembali.", "Periksa body permintaan")
	}

	// Retrieve existing KodeReferral values
	existingKodeReferral := certificate["data"].(bson.M)["kode_referral"].(bson.M)

	// Calculate total JP and total score for hard skills
	totalHSJP, totalHSSkor := uint64(0), float64(0)
	for _, hs := range input.Data.HardSkills.Skills {
		totalHSJP += uint64(hs.SkillJP)
		totalHSSkor += hs.SkillScore
	}

	// Calculate total JP and total score for soft skills
	totalSSJP, totalSSSkor := uint64(0), float64(0)
	for _, ss := range input.Data.SoftSkills.Skills {
		totalSSJP += uint64(ss.SkillJP)
		totalSSSkor += ss.SkillScore
	}

	// Calculate final score
	totalJP := totalHSJP + totalSSJP
	totalScore := totalHSSkor + totalSSSkor
	finalSkor := float64(0)
	if totalJP > 0 {
		finalSkor = float64(math.Round((totalScore/float64(totalJP))*10) / 10)
	}

	// Debugging: Log calculated values
	fmt.Printf("Total HS JP: %d, Total HS Score: %.2f\n", totalHSJP, totalHSSkor)
	fmt.Printf("Total SS JP: %d, Total SS Score: %.2f\n", totalSSJP, totalSSSkor)
	fmt.Printf("Total JP: %d, Total Score: %.2f, Final Skor: %.2f\n", totalJP, totalScore, finalSkor)

	// Update the certificate with the new data while preserving KodeReferral fields
	update := bson.M{
		"$set": bson.M{
			"data": bson.M{
				"sertif_name": input.Data.SertifName,
				"logo":        input.Data.Logo,
				"kode_referral": bson.M{
					"referral_id": existingKodeReferral["referral_id"], // Preserve existing referral ID
					"divisi":      input.Data.KodeReferral.Divisi,      // Update only Divisi
					"bulan_rilis": existingKodeReferral["bulan_rilis"], // Preserve existing Bulan Rilis
					"tahun_rilis": existingKodeReferral["tahun_rilis"], // Preserve existing Tahun Rilis
				},
				"nama_peserta":    input.Data.NamaPeserta,
				"skkni":           input.Data.Skkni,
				"kompeten_bidang": input.Data.KompetenBidang,
				"kompetensi":      input.Data.Kompetensi,
				"validation":      input.Data.Validation,
				"valid_date":      input.Data.ValidDate,
				"total_meet":      input.Data.TotalMeet,
				"meet_time":       input.Data.MeetTime,
				"hard_skills":     input.Data.HardSkills,
				"soft_skills":     input.Data.SoftSkills,
				"final_skor":      float64(math.Round((totalHSSkor+totalSSSkor)/float64(len(input.Data.HardSkills.Skills)+len(input.Data.SoftSkills.Skills))*10) / 10), // Set the calculated final score
			},
			"updated_at": time.Now(),
		},
	}

	_, err = certificateCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal memperbarui Data Sertifikat! Silakan coba lagi.", "Gagal memperbarui Data Sertifikat")
	}

	// Return success
	return OK(c, "Data Sertifikat berhasil diperbarui!", update)
}

// Function for soft delete data certificate
func DeleteCertificate(c *fiber.Ctx) error {
	// Get id param
	idParam := c.Params("id")

	// Convert idParam to ObjectID
	certifID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Sertifikat ini tidak ada! Silakan periksa ID yang dimasukkan.", "Please provide a valid ObjectID")
	}

	// make filter to find document based on acc_id (incremental id)
	filter := bson.M{"_id": certifID}

	// find admin account
	var certificate bson.M
	err = certificateCollection.FindOne(context.TODO(), filter).Decode(&certificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("error: %v\n", err.Error())
			return NotFound(c, "Tidak dapat menemukan sertifikat! Silakan periksa ID yang dimasukkan.", "Cannot find certificate")
		}
		return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", "server error cannot find certificate")
	}

	// Check if the certificate has been deleted
	if deletedAt, ok := certificate["deleted_at"]; ok && deletedAt != nil {
		return AlreadyDeleted(c, "Sertifikat ini telah dihapus! Silakan hubungi admin.", "Check deleted certificate", deletedAt)
	}

	// Update the deleted_at timestamp
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	// update document in collection MongoDB
	result, err := certificateCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal menghapus sertifikat! Silakan coba lagi.", "Delete certificate")
	}

	// Check if the document was already deleted
	if result.MatchedCount == 0 {
		return NotFound(c, "Kompetensi ini tidak dapat ditemukan! Silakan periksa ID yang dimasukkan.", "Found certificate")
	}

	// Respons success
	return OK(c, "Berhasil menghapus sertifikat!", idParam)
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
		return Conflict(c, "Tidak dapat mengunduh sertifikat! Silahkan coba lagi.", err.Error())
	}
	if pdfBytes, err := json.Marshal(resp["data"]); err != nil {
		return Conflict(c, "Tidak dapat mengunduh sertifikat! Silahkan coba lagi.", err.Error())
	} else {
		if err := json.Unmarshal(pdfBytes, &pdf); err != nil {
			return Conflict(c, "Tidak dapat mengunduh sertifikat! Silahkan coba lagi.", err.Error())
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
					return Conflict(c, "Tidak dapat mengunduh sertifikat! Silahkan coba lagi.", err.Error())
				}
			}
		} else {
			return Conflict(c, "Tidak dapat mengunduh sertifikat! Silahkan coba lagi.", err.Error())
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
