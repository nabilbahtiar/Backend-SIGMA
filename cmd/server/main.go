package main

import (
	"fmt"

	"server-room-auth/internal/config"
	"server-room-auth/internal/handler"
	"server-room-auth/internal/model"
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
	
	// AutoMigrate untuk memastikan struktur tabel terbaru
	database.DB.AutoMigrate(
		&model.User{},
		&model.Sensor{},
		&model.SensorConfiguration{},
		&model.SensorTelemetry{},
		&model.Alarm{},
		&model.Notification{},
		&model.HelpdeskTicket{},
		&model.SystemState{},
		&model.AuditLog{},
	)
	
	// Jalankan Seeder
	database.SeedUsers()
	database.SeedSensors()

	// 3. Init Dependency Injection
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	sensorRepo := repository.NewSensorRepository(database.DB)
	sensorService := service.NewSensorService(sensorRepo)
	sensorHandler := handler.NewSensorHandler(sensorService)

	// 4. Setup Router
	r := router.SetupRouter(authHandler, sensorHandler)

	// 5. Jalankan Server
	port := config.AppConfig.AppPort
	fmt.Printf("Server IoT Backend berjalan di http://localhost:%s\n", port)
	r.Run(":" + port)
}
