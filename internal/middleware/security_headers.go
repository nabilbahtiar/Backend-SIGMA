package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders menambahkan HTTP Security Headers standar industri
// untuk mencegah berbagai serangan web (XSS, Clickjacking, MIME Sniffing, dll)
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mencegah Clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Mencegah MIME Type Sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Mengaktifkan XSS Filter browser
		c.Header("X-XSS-Protection", "1; mode=block")
		// Memaksa HTTPS (aktifkan saat sudah pakai TLS/SSL di production)
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Batasi sumber konten yang boleh dimuat (CSP)
		c.Header("Content-Security-Policy", "default-src 'self'")
		// Sembunyikan informasi server
		c.Header("Server", "IoT-Monitor")
		// Kontrol referrer yang dikirim
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}
