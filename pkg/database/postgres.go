package database

import (
	"fmt"
	"log"

	"server-room-auth/internal/config"
	"server-room-auth/internal/model"
	"server-room-auth/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	DB.AutoMigrate(&model.User{})
	fmt.Println("Koneksi database berhasil!")

	SeedUsers()
}

func SeedUsers() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count == 0 {
		fmt.Println("Database kosong. Memasukkan 9 Role awal otomatis (Seeding)...")

		hashedPassword, err := utils.HashPassword("password123")
		if err != nil {
			log.Fatalf("Gagal melakukan hashing password: %v", err)
		}

		users := []model.User{
			{Username: "superadmin", PasswordHash: hashedPassword, Role: "Super Admin IT"},
			{Username: "infra_admin", PasswordHash: hashedPassword, Role: "IT Infrastructure Admin"},
			{Username: "it_support", PasswordHash: hashedPassword, Role: "IT Support"},
			{Username: "net_admin", PasswordHash: hashedPassword, Role: "Network Admin"},
			{Username: "facility_eng", PasswordHash: hashedPassword, Role: "Facility/Engineering"},
			{Username: "security_jaga", PasswordHash: hashedPassword, Role: "Security/Petugas Jaga"},
			{Username: "manajemen", PasswordHash: hashedPassword, Role: "Manajemen"},
			{Username: "auditor", PasswordHash: hashedPassword, Role: "Auditor/Internal Control"},
			{Username: "viewer", PasswordHash: hashedPassword, Role: "Viewer"},
		}

		DB.Create(&users)
		fmt.Println("Data pengguna awal berhasil dibuat!")
	}
}
