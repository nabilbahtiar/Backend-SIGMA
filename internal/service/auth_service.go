package service

import (
	"errors"

	"server-room-auth/internal/repository"
	"server-room-auth/pkg/jwt"
	"server-room-auth/pkg/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

func (s *AuthService) Login(username, password string) (string, string, error) {
	// 1. Cek username
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", "", errors.New("Username yang anda masukkan tidak terdaftar")
	}

	// 2. Cek password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", errors.New("Password yang anda masukkan salah")
	}

	// 3. Generate Token
	token, err := jwt.GenerateJWT(user.Username, user.Role)
	if err != nil {
		return "", "", errors.New("Gagal menghasilkan token sesi")
	}

	return token, user.Role, nil
}
