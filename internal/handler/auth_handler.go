package handler

import (
	"net/http"

	"server-room-auth/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

type LoginRequest struct {
	NIK      string `json:"nik" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid (nik dan password wajib diisi)"})
		return
	}

	// Simpan ke context agar Logger atau Rate Limiter punya info tambahan
	c.Set("attempted_nik", req.NIK)

	token, role, err := h.authService.Login(req.NIK, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"nik":  req.NIK,
			"role": role,
		},
	})
}
