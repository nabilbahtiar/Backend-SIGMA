package middleware

import (
	"log"
	"time"

	"server-room-auth/internal/model"
	"server-room-auth/pkg/database"

	"github.com/gin-gonic/gin"
)

// AuditLogger mencatat setiap request API ke database dan terminal
func AuditLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Proses request terlebih dahulu
		c.Next()

		// Kumpulkan data setelah request selesai diproses
		latency := time.Since(startTime).String()
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		if userAgent == "" {
			userAgent = "Unknown Client"
		}

		// Cari tahu siapa pelakunya (Username)
		username := "Guest"
		if val, exists := c.Get("username"); exists {
			username = val.(string) // Didapat dari JWT jika akses rute terproteksi
		} else if val, exists := c.Get("attempted_username"); exists {
			username = val.(string) // Didapat dari body JSON jika mencoba login
		}

		// Tentukan aksi (Action/Deskripsi singkat)
		action := "API Access"
		if path == "/api/login" {
			if statusCode == 200 {
				action = "Login Success"
			} else {
				action = "Login Failed"
			}
		} else if statusCode == 403 {
			action = "Access Denied (RBAC Block)"
		}

		// 1. Catat ke Terminal
		log.Printf("[AUDIT] %s | %3d | %10s | %-6s | %s | User: %s",
			time.Now().Format("15:04:05"),
			statusCode,
			latency,
			method,
			path,
			username,
		)
		if action == "Login Failed" {
			log.Printf("[SECURITY WARNING] Login gagal dari IP: %s (User: %s)", clientIP, username)
		}

		// 2. Simpan ke Database
		audit := model.AuditLog{
			Timestamp:  startTime,
			Username:   username,
			Action:     action,
			Method:     method,
			Path:       path,
			StatusCode: statusCode,
			ClientIP:   clientIP,
			UserAgent:  userAgent,
			Latency:    latency,
		}
		
		// Gunakan go-routine agar penyimpanan ke database tidak membuat response API menjadi lambat
		go func(logData model.AuditLog) {
			database.DB.Create(&logData)
		}(audit)
	}
}
