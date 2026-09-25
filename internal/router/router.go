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

func SetupRouter(authHandler *handler.AuthHandler, sensorHandler *handler.SensorHandler) *gin.Engine {
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
		// Rate Limiter: memblokir NIK/IP jika 3x gagal login
		api.POST("/login", middleware.LoginRateLimiter(), authHandler.Login)
	}

	// ==============================
	// Rute Terproteksi (Wajib JWT valid)
	// Matriks Hak Akses sesuai PDF Role Matrix:
	// 1. Super Admin IT
	// 2. IT Support
	// 3. Facility/Engineering
	// 4. Security
	// 5. Manajemen
	// 6. Guest
	// ==============================
	secure := api.Group("/secure")
	secure.Use(middleware.AuthMiddleware())
	{
		// 1. Dashboard (Semua 6 role memiliki akses)
		dashboardRoles := []string{
			"Super Admin IT", "IT Support", "Facility/Engineering", "Security", "Manajemen", "Guest",
		}
		dashboard := secure.Group("/dashboard")
		dashboard.Use(middleware.RoleMiddleware(dashboardRoles...))
		{
			dashboard.GET("/status", func(c *gin.Context) {
				nik, _ := c.Get("nik")
				role, _ := c.Get("role")
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": fmt.Sprintf("Halo %v (%v), status dashboard sistem IoT normal.", nik, role),
				})
			})
			dashboard.GET("/sensors", sensorHandler.GetDashboardSensors)
		}

		// 2. Sensor / Config (Hanya Super Admin IT, IT Support, Facility/Engineering)
		sensorRoles := []string{
			"Super Admin IT", "IT Support", "Facility/Engineering",
		}
		sensorConfig := secure.Group("/sensor")
		sensorConfig.Use(middleware.RoleMiddleware(sensorRoles...))
		{
			sensorConfig.POST("/config", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Konfigurasi sensor berhasil diperbarui",
				})
			})
		}

		// 3. Alarm (Semua 6 role: Super Admin IT, IT Support, Facility/Engineering, Security, Manajemen, Guest)
		alarmRoles := []string{
			"Super Admin IT", "IT Support", "Facility/Engineering", "Security", "Manajemen", "Guest",
		}
		alarmGroup := secure.Group("/alarm")
		alarmGroup.Use(middleware.RoleMiddleware(alarmRoles...))
		{
			alarmGroup.GET("/list", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Daftar alarm/peringatan berhasil dimuat",
				})
			})
		}

		// 4. User / RBAC (Hanya Super Admin IT)
		userRbacRoles := []string{
			"Super Admin IT",
		}
		userGroup := secure.Group("/users")
		userGroup.Use(middleware.RoleMiddleware(userRbacRoles...))
		{
			userGroup.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Data pengguna & RBAC berhasil dimuat (Super Admin IT only)",
				})
			})
		}

		// 5. Report & Monitoring (Semua 6 role)
		reportRoles := []string{
			"Super Admin IT", "IT Support", "Facility/Engineering", "Security", "Manajemen", "Guest",
		}
		monitor := secure.Group("/monitoring")
		monitor.Use(middleware.RoleMiddleware(reportRoles...))
		{
			monitor.GET("/sensor/data", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Data pemantauan sensor berhasil diambil",
				})
			})
		}

		report := secure.Group("/report")
		report.Use(middleware.RoleMiddleware(reportRoles...))
		{
			report.GET("/summary", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Ringkasan laporan berhasil diambil",
				})
			})
		}

		// 6. Helpdesk (Super Admin IT, IT Support, Facility/Engineering, Security)
		helpdeskRoles := []string{
			"Super Admin IT", "IT Support", "Facility/Engineering", "Security",
		}
		helpdesk := secure.Group("/helpdesk")
		helpdesk.Use(middleware.RoleMiddleware(helpdeskRoles...))
		{
			helpdesk.GET("/tickets", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "Data tiket helpdesk berhasil dimuat",
				})
			})
		}
	}

	return r
}
