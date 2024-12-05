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

	pdfReq.Data = *processCertificate(&pdfReq.Data)

	mappedData := model.CertificateData{
		AdminId:     objectID,
		SertifName:  pdfReq.Data.SertifName,
		SertifTitle: pdfReq.Data.SertifTitle,
		KodeReferral: model.KodeReferral{
			ReferralID: nextReferralID,
			Divisi:     kompetensi.Divisi,
			BulanRilis: monthRoman,
			TahunRilis: year,
		},
		NamaPeserta:    pdfReq.Data.NamaPeserta,
		SKKNI:          kompetensi.SKKNI,
		KompetenBidang: pdfReq.Data.KompetenBidang,
		Kompetensi:     pdfReq.Data.Kompetensi,
		Validation:     pdfReq.Data.Validation,
		DataID:         newDataID,
		TotalJP:        pdfReq.Data.TotalJP,
		TotalMeet:      pdfReq.Data.TotalMeet,
		MeetTime:       pdfReq.Data.MeetTime,
		ValidDate:      pdfReq.Data.ValidDate,
		HardSkills:     pdfReq.Data.HardSkills,
		SoftSkills:     pdfReq.Data.SoftSkills,
		FinalSkor:      pdfReq.Data.FinalSkor,
		Signature: model.SignatureData{
			ConfigName: pdfReq.Data.Signature.ConfigName,
			Stamp:      pdfReq.Data.Signature.Stamp,
			Signature:  pdfReq.Data.Signature.Signature,
			Logo:       pdfReq.Data.Signature.Logo,
			Name:       pdfReq.Data.Signature.Name,
			Role:       pdfReq.Data.Signature.Role,
		},
	}

	certificate := model.PDF{
		AdminId:     objectID,
		DataID:      newDataID,
		SertifName:  pdfReq.Data.SertifName,
		SertifTitle: pdfReq.Data.SertifTitle,
		Data:        mappedData,
		Model: model.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
	}

	// make pdf creation concurrent to return handler faster
	go generator.CreatePDF(c, &mappedData, "ab", true)

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
		"_id":          1,
		"admin_id":     1,
		"data_id":      1,
		"sertif_name":  1,
		"sertif_title": 1,
		"created_at":   1,
		"updated_at":   1,
		"deleted_at":   1,
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
	oID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Data Sertifikat tidak ditemukan! Silakan periksa ID sertifikat.", "Gagal mengonversi parameter pada Edit Sertifikat")
	}

	filter := bson.M{"_id": oID}
	var certificateMongo model.PDF
	if err := certificateCollection.FindOne(c.Context(), filter).Decode(&certificateMongo); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Data Sertifikat tidak ditemukan!", "Gagal menemukan Data Sertifikat")
		}
		return Conflict(c, "Gagal mendapatkan Data Sertifikat! Silakan coba lagi.", "Gagal menemukan Data Sertifikat")
	}

	if certificateMongo.DeletedAt != nil {
		return AlreadyDeleted(c, "Data Sertifikat ini sudah dihapus! Silakan hubungi Data Sertifikat.", "Periksa Data Sertifikat yang dihapus", certificateMongo.DeletedAt)
	}

	var pdf model.PDF
	if err := c.BodyParser(&pdf); err != nil {
		return BadRequest(c, "Gagal mendapatkan Data Sertifikat! Silakan coba lagi.", err.Error())
	}
	pdfData := pdf.Data

	// fetch Kompetensi by the given nama_kompetensi from the request
	var kompetensi model.Kompetensi
	if err := competenceCollection.FindOne(context.TODO(), bson.M{"nama_kompetensi": pdfData.Kompetensi}).Decode(&kompetensi); err != nil {
		return NotFound(c, "Tidak menemukan kompetensi yang dicari. Mohon coba lagi.", err.Error())
	}

	// Retrieve existing KodeReferral values
	if refBytes, err := json.Marshal(certificateMongo); err != nil {
		return BadRequest(c, "Gagal mendapatkan Data Sertifikat! Silakan coba lagi.", err.Error())
	} else {
		if err := json.Unmarshal(refBytes, &pdfData.KodeReferral); err != nil {
			return BadRequest(c, "Gagal mendapatkan Data Sertifikat! Silakan coba lagi.", err.Error())
		}
	}

	pdfData = *processCertificate(&pdfData)

	mappedData := model.CertificateData{
		AdminId:     certificateMongo.AdminId,
		SertifName:  pdfData.SertifName,
		SertifTitle: pdfData.SertifTitle,
		KodeReferral: model.KodeReferral{
			ReferralID: certificateMongo.Data.KodeReferral.ReferralID,
			Divisi:     kompetensi.Divisi,
			BulanRilis: certificateMongo.Data.KodeReferral.BulanRilis,
			TahunRilis: certificateMongo.Data.KodeReferral.TahunRilis,
		},
		NamaPeserta:    pdfData.NamaPeserta,
		SKKNI:          kompetensi.SKKNI,
		KompetenBidang: pdfData.KompetenBidang,
		Kompetensi:     pdfData.Kompetensi,
		Validation:     pdfData.Validation,
		DataID:         certificateMongo.DataID,
		TotalJP:        pdfData.TotalJP,
		TotalMeet:      pdfData.TotalMeet,
		MeetTime:       pdfData.MeetTime,
		ValidDate:      pdfData.ValidDate,
		HardSkills:     pdfData.HardSkills,
		SoftSkills:     pdfData.SoftSkills,
		FinalSkor:      pdfData.FinalSkor,
		Signature: model.SignatureData{
			ConfigName: pdfData.Signature.ConfigName,
			Stamp:      pdfData.Signature.Stamp,
			Signature:  pdfData.Signature.Signature,
			Logo:       pdfData.Signature.Logo,
			Name:       pdfData.Signature.Name,
			Role:       pdfData.Signature.Role,
		},
	}

	certificate := model.PDF{
		AdminId:     certificateMongo.AdminId,
		DataID:      certificateMongo.DataID,
		SertifName:  pdfData.SertifName,
		SertifTitle: pdfData.SertifTitle,
		Data:        mappedData,
		Model: model.Model{
			ID:        certificateMongo.ID,
			CreatedAt: certificateMongo.CreatedAt,
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// Update the certificate with the new data while preserving KodeReferral fields
	update := bson.M{"$set": certificate}

	go generator.CreatePDF(c, &pdfData, "ab", true)

	if _, err := certificateCollection.UpdateOne(c.Context(), filter, update); err != nil {
		return Conflict(c, "Gagal memperbarui Data Sertifikat! Silakan coba lagi.", err.Error())
	}

	// Return success
	return OK(c, "Data Sertifikat berhasil diperbarui!", certificate)
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
			if err = generator.CreatePDF(c, &data, certifType, false); err != nil {
				return Conflict(c, "Tidak dapat mengunduh sertifikat! Silahkan coba lagi.", err.Error())
			}
		} else {
			return Conflict(c, "Tidak dapat mengunduh sertifikat! Silahkan coba lagi.", err.Error())
		}
	}

	c.Response().Header.Add("Content-Type", "application/pdf")
	addAllowOrigin(c)

	return c.Download(filepath, "Sertifikat BTW Edutech "+certifType+" - "+data.NamaPeserta)
}

func processCertificate(certif *model.CertificateData) *model.CertificateData {
	certif.SertifName = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.ToUpper(certif.SertifName)), "SERTIFIKAT"))
	certif.KodeReferral.Divisi = strings.ToUpper(certif.KodeReferral.Divisi)
	certif.SertifTitle = fmt.Sprintf("%s - %s", certif.NamaPeserta, certif.Kompetensi)

	totalHSJP, totalHSSkor := uint64(0), float64(0)
	for _, hs := range certif.HardSkills.Skills {
		totalHSJP += hs.SkillJP
		totalHSSkor += hs.SkillScore
	}

	totalSSJP, totalSSSkor := uint64(0), float64(0)
	for _, ss := range certif.SoftSkills.Skills {
		totalSSJP += ss.SkillJP
		totalSSSkor += ss.SkillScore
	}

	certif.TotalJP = totalHSJP + totalSSJP
	certif.FinalSkor = float64(math.Round((totalHSSkor+totalSSSkor)/float64(len(certif.HardSkills.Skills)+len(certif.SoftSkills.Skills))*10) / 10)

	certif.HardSkills = model.SkillPDF{
		Skills:          certif.HardSkills.Skills,
		TotalSkillJP:    totalHSJP,
		TotalSkillScore: float64(math.Round(totalHSSkor/float64(len(certif.HardSkills.Skills))*10) / 10),
	}

	certif.SoftSkills = model.SkillPDF{
		Skills:          certif.SoftSkills.Skills,
		TotalSkillJP:    totalSSJP,
		TotalSkillScore: float64(math.Round(totalSSSkor/float64(len(certif.SoftSkills.Skills))*10) / 10),
	}

	return certif
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
