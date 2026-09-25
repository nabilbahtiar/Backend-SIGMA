package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type IPBanState struct {
	FailCount   int
	BanCount    int
	BannedUntil time.Time
}

var (
	// ipBans sekarang menggunakan kombinasi "IP_NIK" sebagai key
	ipBans = make(map[string]*IPBanState)
	banMu  sync.Mutex
)

// LoginRateLimiter memblokir IP yang gagal login 3x dengan skema waktu bertingkat (5m, 10m, 20m, dst)
func LoginRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		banMu.Lock()
		state, exists := ipBans[ip]
		if !exists {
			state = &IPBanState{}
			ipBans[ip] = state
		}

		// 2. Cek apakah IP ini sedang dihukum
		if time.Now().Before(state.BannedUntil) {
			remainingSeconds := int(time.Until(state.BannedUntil).Seconds())
			banMu.Unlock()
			if remainingSeconds < 1 {
				remainingSeconds = 1
			}
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":               "Terlalu banyak percobaan gagal dari IP ini. Silakan coba lagi nanti.",
				"retry_after_seconds": remainingSeconds,
			})
			c.Abort()
			return
		}
		banMu.Unlock()

		// 3. Lanjutkan request ke Auth Handler
		c.Next()

		// 4. Evaluasi hasil request
		status := c.Writer.Status()

		banMu.Lock()
		defer banMu.Unlock()

		switch status {
		case http.StatusUnauthorized, http.StatusBadRequest:
			// Login gagal
			state.FailCount++

			if state.FailCount >= 3 {
				state.BanCount++ // Tingkatkan level hukuman

				// Hitung masa hukuman (5 menit * 2^(BanCount-1))
				multiplier := 1 << uint(state.BanCount-1)
				banDuration := time.Duration(5*multiplier) * time.Minute

				state.BannedUntil = time.Now().Add(banDuration)
				state.FailCount = 0 // Reset hitungan kegagalan untuk masa depan
			}
		case http.StatusOK:
			// Jika login berhasil, hapus semua catatan buruk untuk IP ini
			delete(ipBans, ip)
		}
	}
}
