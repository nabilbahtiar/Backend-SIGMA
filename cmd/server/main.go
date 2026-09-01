package main

import (
	"fmt"

	"server-room-auth/internal/config"
	"server-room-auth/internal/handler"
	"server-room-auth/internal/repository"
	"server-room-auth/internal/router"
	"server-room-auth/internal/service"
	"server-room-auth/pkg/database"
	"server-room-auth/pkg/jwt"
)

func main() {
	// 1. Load Konfigurasi
	config.LoadConfig()

	// 2. Init Utilitas & Database
	jwt.InitJWTSecret(config.AppConfig.JWTSecret)
	database.InitDB()

	// 3. Init Dependency Injection
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// 4. Setup Router
	r := router.SetupRouter(authHandler)

	// 5. Jalankan Server
	port := config.AppConfig.AppPort
	fmt.Printf("Server IoT Backend berjalan di http://localhost:%s\n", port)
	r.Run(":" + port)
}
