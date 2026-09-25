package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JWTSecret []byte

const (
	TokenIssuer   = "iot-monitor-rs-citra-husada"
	TokenDuration = 24 * time.Hour // Token berlaku 24 jam
)

// InitJWTSecret menginisialisasi secret dari environment
func InitJWTSecret(secret string) {
	JWTSecret = []byte(secret)
}

// GenerateJWT menghasilkan token JWT yang aman dengan claims standar
func GenerateJWT(nik, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		// Standard Claims
		"iss": TokenIssuer,                          // Issuer: siapa yang mengeluarkan token
		"jti": uuid.New().String(),                  // JWT ID: unik untuk setiap token (mencegah replay attack)
		"iat": now.Unix(),                           // Issued At: waktu token dibuat
		"exp": now.Add(TokenDuration).Unix(),        // Expiry: kapan token kadaluarsa
		"nbf": now.Unix(),                           // Not Before: token tidak valid sebelum waktu ini

		// Custom Claims (data aplikasi)
		"nik":  nik,
		"role": role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ParseAndVerify memeriksa apakah token JWT valid dan memverifikasi issuer
func ParseAndVerify(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode algoritma (WAJIB, mencegah 'none' algorithm attack)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode penandatanganan tidak valid: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	},
		// Validasi issuer (mencegah token dari sistem lain diterima)
		jwt.WithIssuer(TokenIssuer),
		// Validasi waktu expired secara otomatis
		jwt.WithExpirationRequired(),
	)

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token tidak valid atau sudah kadaluarsa")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("gagal mengekstrak payload token")
	}

	return claims, nil
}
