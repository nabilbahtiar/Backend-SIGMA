package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLogger mencatat setiap request API ke konsol (siap dialihkan ke file/DB)
func AuditLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Proses request
		c.Next()

		// Catat setelah request selesai
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[AUDIT] %s | %3d | %13v | %-7s | %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			method,
			path,
		)

		// Catat percobaan login gagal sebagai peringatan
		if path == "/api/login" && statusCode == 401 {
			log.Printf("[SECURITY WARNING] Login gagal dari IP: %s | Path: %s", clientIP, path)
		}
	}
}
