package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ==============================
// 1. Model & Konfigurasi
// ==============================

// User merepresentasikan tabel users di database
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"` // Tidak dikirim dalam JSON
	Role         string    `gorm:"not null" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var DB *gorm.DB
var JWTSecret = []byte("Rahasia_Super_Aman_Sistem_IoT_123!") // Ganti via Environment Variable di Production

// InitDB menghubungkan Golang ke PostgreSQL menggunakan GORM
func InitDB() {
	// Sesuaikan konfigurasi ini dengan PostgreSQL lokal Anda
	dsn := "host=localhost user=postgres password=123 dbname=iot_server_room port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Auto-Migrate tabel (opsional, jika Anda belum menjalankan SQL di atas)
	DB.AutoMigrate(&User{})
	fmt.Println("Koneksi database berhasil!")

	// Jalankan seeder
	SeedUsers()
}

// SeedUsers berfungsi untuk mengisi data user awal jika tabel users masih kosong
func SeedUsers() {
	var count int64
	DB.Model(&User{}).Count(&count)
	if count == 0 {
		fmt.Println("Menyiapkan data awal pengguna (Seeding)...")
		
		// Password default untuk semua user: password123
		hashedPassword, err := HashPassword("password123")
		if err != nil {
			log.Fatalf("Gagal melakukan hashing password: %v", err)
		}
		
		users := []User{
			{Username: "superadmin", PasswordHash: hashedPassword, Role: "Super Admin"},
			{Username: "infra_admin", PasswordHash: hashedPassword, Role: "IT Infrastructure Admin"},
			{Username: "it_support", PasswordHash: hashedPassword, Role: "IT Support"},
			{Username: "net_admin", PasswordHash: hashedPassword, Role: "Network Admin"},
			{Username: "security_jaga", PasswordHash: hashedPassword, Role: "Security"},
		}
		
		DB.Create(&users)
		fmt.Println("Data pengguna awal berhasil dibuat!")
	}
}

// ==============================
// 2. Utilitas Keamanan (Bcrypt & JWT)
// ==============================

// HashPassword mengenkripsi password (digunakan saat pendaftaran/pembuatan user)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash membandingkan teks asli dengan hash di database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT menghasilkan token JWT berdasarkan username dan role
func GenerateJWT(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ==============================
// 3. API Handlers (Login)
// ==============================

// LoginRequest adalah payload JSON yang diharapkan dari frontend
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler menangani proses autentikasi
func LoginHandler(c *gin.Context) {
	var req LoginRequest

	// Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid (username dan password wajib diisi)"})
		return
	}

	// Cari user berdasarkan username di database
	var user User
	if err := DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "username yang anda masukkan tidak terdaftar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan pada server"})
		return
	}

	// Validasi kecocokan password
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password yang anda masukkan salah"})
		return
	}

	// Buat JWT Token
	token, err := GenerateJWT(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan token sesi"})
		return
	}

	// Berikan respon berhasil
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// ==============================
// 4. Middleware (Autentikasi & RBAC)
// ==============================

// AuthMiddleware memvalidasi token JWT dari header Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Authorization tidak ditemukan"})
			c.Abort()
			return
		}

		// Format token yang diharapkan: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format Authorization header tidak valid (gunakan Bearer)"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi metode algoritma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Metode penandatanganan tidak valid: %v", token.Header["alg"])
			}
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa"})
			c.Abort()
			return
		}

		// Ambil claims dari token dan simpan di context agar bisa dipakai oleh handler
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
		}

		c.Next()
	}
}

// RoleMiddleware (RBAC) membatasi akses berdasarkan array role yang diizinkan
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: role tidak terdefinisi"})
			c.Abort()
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: hak akses tidak mencukupi (Privilege Escalation Blocked)"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ==============================
// 5. Main (Routing & Eksekusi)
// ==============================

func main() {
	InitDB()

	r := gin.Default()

	// Setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Izinkan semua origin saat tahap development
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Route publik
	r.POST("/api/login", LoginHandler)

	// Route terproteksi (Membutuhkan JWT valid)
	protected := r.Group("/api/secure")
	protected.Use(AuthMiddleware())
	{
		// Semua user terotentikasi bisa mengakses ini
		protected.GET("/dashboard/status", func(c *gin.Context) {
			username, _ := c.Get("username")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Halo %v (%v), sistem IoT normal.", username, role),
			})
		})

		// RBAC: Hanya Super Admin & IT Infrastructure Admin yang bisa mengatur ambang batas sensor
		adminOnly := protected.Group("/sensor")
		adminOnly.Use(RoleMiddleware("Super Admin", "IT Infrastructure Admin"))
		{
			adminOnly.POST("/config", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Konfigurasi sensor berhasil diperbarui"})
			})
		}
	}

	fmt.Println("Server IoT Backend berjalan di http://localhost:8080")
	r.Run(":8080")
}
