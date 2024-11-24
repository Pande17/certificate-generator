package rest

import (
	"context"
	"os"
	"time"

	"certificate-generator/database"
	"certificate-generator/model"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Function to make new admin account
func SignUp(c *fiber.Ctx) error {
	var adminReq struct {
		AdminName     string `json:"admin_name" valid:"required~Nama tidak boleh kosong!, stringlength(1|30)~Nama harus antara 1 hingga 30 karakter!"`
		AdminPassword string `json:"admin_password" valid:"required~Password tidak boleh kosong!"`
	}

	// Parse the request body
	if err := c.BodyParser(&adminReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", err.Error())
	}

	// Validate the input data using govalidator
	if _, err := govalidator.ValidateStruct(&adminReq); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", err.Error())
	}

	// Hashing the input password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(adminReq.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return Conflict(c, "Gagal membuat akun! Silakan coba lagi.", "Gagal melakukan hashing password")
	}

	// Connect to the admin account collection in the database
	adminCollection := database.GetCollection("adminAcc")

	// Check the availability of the admin account name
	var existingAdmin model.AdminAccount
	filter := bson.M{"admin_name": adminReq.AdminName}

	// Find admin account with the same account name as input name
	err = adminCollection.FindOne(context.TODO(), filter).Decode(&existingAdmin)
	if err == nil {
		return BadRequest(c, "Nama ini sudah digunakan! Silakan pilih nama lain.", "Periksa nama admin")
	} else if err != mongo.ErrNoDocuments {
		return Conflict(c, "Kesalahan server! Silakan coba lagi.", "Periksa nama admin")
	}

	// Input data from req struct to struct "AdminAccount"
	admin := model.AdminAccount{
		AdminName:     adminReq.AdminName,
		AdminPassword: string(hash),
		Model: model.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// Insert data from struct "AdminAccount" to collection in database MongoDB
	_, err = adminCollection.InsertOne(context.TODO(), admin)
	if err != nil {
		return Conflict(c, "Gagal membuat akun admin! Silakan coba lagi.", "Gagal menyimpan akun admin")
	}

	// Return success
	return OK(c, "Akun Admin berhasil dibuat!", admin)
}

// Function to login to admin account
func Login(c *fiber.Ctx) error {
	var adminReq struct {
		AdminName     string `json:"admin_name" valid:"required~Nama tidak boleh kosong!, stringlength(1|30)~Nama harus antara 1 hingga 30 karakter!"`
		AdminPassword string `json:"admin_password" valid:"required~Password tidak boleh kosong!"`
	}

	// Parse the request body
	if err := c.BodyParser(&adminReq); err != nil {
		return BadRequest(c, "Input tidak valid! Silakan periksa kembali.", "Gagal mem-parsing body login")
	}

	// New variable to store admin login data
	var admin model.AdminAccount
	adminCollection := database.GetCollection("adminAcc")

	// Find admin account with the same account name as input name
	err := adminCollection.FindOne(context.TODO(), bson.M{"admin_name": adminReq.AdminName}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Unauthorized(c, "Nama atau Password salah! Silakan coba lagi.", "Gagal menemukan nama admin saat login")
		}
		return Conflict(c, "Kesalahan database! Silakan coba lagi.", "Gagal menemukan nama admin saat login")
	}

	// Check if DeletedAt field already has a value (account has been deleted)
	if admin.DeletedAt != nil && !admin.DeletedAt.IsZero() {
		return AlreadyDeleted(c, "Akun ini telah dihapus! Silakan hubungi admin.", "Periksa akun admin yang dihapus", admin.DeletedAt)
	}

	// Hashing the input password with bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(adminReq.AdminPassword))
	if err != nil {
		return Unauthorized(c, "Nama atau Password salah! Silakan coba lagi.", "Gagal memverifikasi password saat login")
	}

	// Generate token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Retrieve the "SECRET" environment variable
	secret := os.Getenv("SECRET")
	if secret == "" {
		return Conflict(c, "Kunci rahasia tidak diset! Silakan hubungi admin.", "Gagal mengambil kunci rahasia")
	}

	// Use the secret key to sign the token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return BadRequest(c, "Login gagal! Silakan coba lagi.", "Gagal menggunakan kunci rahasia")
	}

	// Set a cookie for admin
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	// Return success
	return OK(c, "Login berhasil! Selamat datang.", tokenString)
}

// Function to validate if the user has a valid authentication cookie
func Validate(c *fiber.Ctx) error {
	adminID := c.Locals("admin")

	if adminID == nil {
		return Unauthorized(c, "Silakan login terlebih dahulu!", "Gagal memvalidasi pengguna")
	}

	// Return a success message along with the admin ID
	return OK(c, "User terverifikasi! Anda sudah login.", adminID)
}

// Function to logout
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})
	// Return success
	return OK(c, "Anda berhasil logout!", nil)
}

// Function to edit admin account data
func EditAdminAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	accID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Akun tidak ditemukan! Silakan periksa ID akun.", "Gagal mengonversi parameter pada Edit Admin")
	}

	adminCollection := database.GetCollection("adminAcc")
	filter := bson.M{"_id": accID}
	var acc bson.M

	if err := adminCollection.FindOne(c.Context(), filter).Decode(&acc); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Akun tidak ditemukan!", "Gagal menemukan akun")
		}
		return Conflict(c, "Gagal mendapatkan akun! Silakan coba lagi.", "Gagal menemukan akun")
	}

	if deletedAt, ok := acc["deleted_at"]; ok && deletedAt != nil {
		return AlreadyDeleted(c, "Akun ini sudah dihapus! Silakan hubungi admin.", "Periksa akun admin yang dihapus", deletedAt)
	}

	var input struct {
		AdminName     string `json:"admin_name" valid:"required~Nama tidak boleh kosong!, stringlength(1|30)~Nama harus antara 1 hingga 30 karakter!"`
		AdminPassword string `json:"admin_password" valid:"required~Password tidak boleh kosong!"`
	}

	if err := c.BodyParser(&input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid! Silakan periksa kembali.", "Periksa body permintaan")
	}

	// Validate the input data using govalidator
	if _, err := govalidator.ValidateStruct(&input); err != nil {
		return BadRequest(c, "Data yang dimasukkan tidak valid!", err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return BadRequest(c, "Gagal memperbarui password! Silakan coba lagi.", "")
	}

	update := bson.M{
		"$set": bson.M{
			"admin_name":     input.AdminName,
			"admin_password": string(hash),
			"updated_at":     time.Now(),
		},
	}

	_, err = adminCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal memperbarui akun admin! Silakan coba lagi.", "Gagal memperbarui akun admin")
	}

	// Return success
	return OK(c, "Akun admin berhasil diperbarui!", update)
}

// Function for soft delete admin account
func DeleteAdminAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	accID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Akun tidak ditemukan! Silakan periksa ID akun.", "Gagal mengonversi parameter")
	}

	adminCollection := database.GetCollection("adminAcc")
	filter := bson.M{"_id": accID}

	var adminAccount bson.M
	err = adminCollection.FindOne(context.TODO(), filter).Decode(&adminAccount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Akun tidak ditemukan!", "Gagal menemukan akun admin")
		}
		return Conflict(c, "Gagal mengambil detail akun! Silakan coba lagi.", "Gagal menemukan akun admin")
	}

	if deletedAt, ok := adminAccount["deleted_at"]; ok && deletedAt != nil {
		return AlreadyDeleted(c, "Akun ini sudah dihapus! Silakan hubungi admin.", "Periksa akun admin yang dihapus", deletedAt)
	}

	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	result, err := adminCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Conflict(c, "Gagal menghapus akun! Silakan coba lagi.", "Gagal menghapus akun admin")
	}

	if result.MatchedCount == 0 {
		return NotFound(c, "Akun tidak ditemukan!", "Gagal menemukan akun admin")
	}

	// Return success
	return OK(c, "Akun berhasil dihapus!", accID)
}

// Function to get admin account
func GetAdminAccount(c *fiber.Ctx) error {
	if len(c.Queries()) == 0 {
		return getAllAdminAccount(c)
	}
	key := c.Query("type")
	val := c.Query("s")
	var value any
	if key == "id" {
		key = "_id"
		var err error
		if value, err = primitive.ObjectIDFromHex(val); err != nil {
			return BadRequest(c, "Tidak dapat mem-parsing id! Silakan periksa kembali.", err.Error())
		}
	} else {
		value = val
	}
	return getOneAdminAccount(c, bson.M{key: value})
}

// Function to see all Admin Accounts
func getAllAdminAccount(c *fiber.Ctx) error {
	var results []bson.M
	adminCollection := database.GetCollection("adminAcc")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projection := bson.M{
		"_id":        1,
		"admin_name": 1,
		"created_at": 1,
		"updated_at": 1,
		"deleted_at": 1,
	}

	cursor, err := adminCollection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Tidak ada akun ditemukan!", err.Error())
		}
		return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", "Gagal menemukan akun admin")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var admin bson.M
		if err := cursor.Decode(&admin); err != nil {
			return Conflict(c, "Gagal mendekode data! Silakan coba lagi.", "Gagal mendekode data")
		}
		results = append(results, admin)
	}
	if err := cursor.Err(); err != nil {
		return Conflict(c, "Kesalahan cursor! Silakan coba lagi.", "Kesalahan cursor")
	}

	return OK(c, "Berhasil mendapatkan semua data akun admin!", results)
}

// Function to get detail account by accID
func getOneAdminAccount(c *fiber.Ctx, filter bson.M) error {
	adminCollection := database.GetCollection("adminAcc")
	var accountDetail bson.M

	if err := adminCollection.FindOne(context.TODO(), filter).Decode(&accountDetail); err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Data tidak ditemukan! Silakan periksa kembali.", "Gagal menemukan akun")
		}
		return Conflict(c, "Gagal mengambil data! Silakan coba lagi.", "Gagal menemukan akun")
	}

	if deletedAt, ok := accountDetail["deleted_at"]; ok && deletedAt != nil {
		return AlreadyDeleted(c, "Akun ini sudah dihapus! Silakan hubungi admin.", "Periksa akun admin yang dihapus", deletedAt)
	}

	// Return success
	return OK(c, "Berhasil mendapatkan data akun!", accountDetail)
}
