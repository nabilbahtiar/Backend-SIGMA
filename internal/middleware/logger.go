package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLogger mencatat setiap request API ke terminal/konsol
func AuditLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Proses request terlebih dahulu
		c.Next()

		// Kumpulkan data setelah request selesai diproses
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		statusText := http.StatusText(statusCode)
		if statusText == "" {
			statusText = "Unknown Status"
		}

		log.Printf("| %3d %-20s | %13v | %-7s | %s",
			statusCode,
			statusText,
			latency,
			method,
			path,
		)

		// Catat percobaan login gagal sebagai peringatan keamanan
		if path == "/api/login" && statusCode == 401 {
			log.Printf("[SECURITY WARNING] Login gagal dari IP: %s | Path: %s", clientIP, path)
		}
	}
}
