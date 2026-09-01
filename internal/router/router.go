package router

import (
	"fmt"
	"net/http"
	"time"

	"server-room-auth/internal/handler"
	"server-room-auth/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler) *gin.Engine {
	// Gunakan gin.New() agar kita kontrol penuh semua middleware
	r := gin.New()

	// ==============================
	// Middleware Global
	// ==============================
	r.Use(gin.Recovery())               // Auto-recover dari panic, server tidak akan mati
	r.Use(middleware.AuditLogger())     // Log setiap request & deteksi login gagal
	r.Use(middleware.SecurityHeaders()) // HTTP Security Headers (XSS, Clickjacking, dll)

	// CORS - Batasi hanya dari origin yang diizinkan
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Izinkan semua origin (development mode)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	// Batasi ukuran body request maks 2MB (anti payload bomb / DoS)
	r.MaxMultipartMemory = 2 << 20

	// ==============================
	// Rute Publik
	// ==============================
	api := r.Group("/api")
	{
		// Rate Limiter: maks 5 percobaan login per menit per IP (anti brute-force)
		api.POST("/login", middleware.LoginRateLimiter(), authHandler.Login)
	}

	// ==============================
	// Rute Terproteksi (Wajib JWT valid)
	// ==============================
	secure := api.Group("/secure")
	secure.Use(middleware.AuthMiddleware())
	{
		// Semua role terautentikasi bisa akses status
		secure.GET("/dashboard/status", func(c *gin.Context) {
			username, _ := c.Get("username")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": fmt.Sprintf("Halo %v (%v), sistem IoT normal.", username, role),
			})
		})

		// RBAC: Hanya Super Admin IT & IT Infrastructure Admin (konfigurasi sensor)
		adminOnly := secure.Group("/sensor")
		adminOnly.Use(middleware.RoleMiddleware("Super Admin IT", "IT Infrastructure Admin"))
		{
			adminOnly.POST("/config", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Konfigurasi sensor berhasil diperbarui",
				})
			})
		}

		// RBAC: Admin + IT Support + Network Admin bisa lihat data monitoring
		monitorRoles := []string{
			"Super Admin IT", "IT Infrastructure Admin",
			"IT Support", "Network Admin",
			"Facility/Engineering", "Security/Petugas Jaga",
		}
		monitor := secure.Group("/monitoring")
		monitor.Use(middleware.RoleMiddleware(monitorRoles...))
		{
			monitor.GET("/sensor/data", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Data sensor berhasil diambil (endpoint placeholder)",
				})
			})
		}

		// RBAC: Manajemen & Auditor hanya bisa lihat laporan (read-only)
		reportRoles := []string{
			"Super Admin IT", "Manajemen",
			"Auditor/Internal Control", "IT Infrastructure Admin",
		}
		report := secure.Group("/report")
		report.Use(middleware.RoleMiddleware(reportRoles...))
		{
			report.GET("/summary", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Laporan berhasil diambil (endpoint placeholder)",
				})
			})
		}
	}

	return r
}
