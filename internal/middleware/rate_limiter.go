package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Menyimpan rate limiter per IP
var (
	loginLimiters = make(map[string]*rate.Limiter)
	mu            sync.Mutex
)

// getLoginLimiter mengembalikan rate limiter untuk IP tertentu
// Batas: 5 request per menit (setiap 12 detik 1 token diisi)
func getLoginLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := loginLimiters[ip]
	if !exists {
		// 5 percobaan per menit (burst 5, rate 1 per 12 detik)
		limiter = rate.NewLimiter(rate.Every(12e9), 5) // 5 req/menit
		loginLimiters[ip] = limiter
	}
	return limiter
}

// LoginRateLimiter mencegah serangan brute force pada endpoint login
func LoginRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLoginLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Terlalu banyak percobaan login. Silakan tunggu beberapa saat.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
