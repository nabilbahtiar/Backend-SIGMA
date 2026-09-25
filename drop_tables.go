package main

import (
	"fmt"
	"log"

	"server-room-auth/internal/config"
	"server-room-auth/internal/model"
	"server-room-auth/pkg/database"
)

func main() {
	config.LoadConfig()
	database.InitDB()

	// Menghapus tabel agar auto-migrate berjalan bersih dari awal
	tables := []string{
		"users", "sensors", "sensor_configurations", "sensor_telemetries",
		"alarms", "notifications", "helpdesk_tickets", "system_states", "audit_logs",
	}

	for _, table := range tables {
		err := database.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", table)).Error
		if err != nil {
			log.Fatalf("Gagal menghapus tabel %s: %v", table, err)
		}
	}

	fmt.Println("Tabel-tabel lama berhasil dibersihkan.")

	// Jalankan AutoMigrate
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
	fmt.Println("AutoMigrate untuk semua model berhasil dijalankan.")

	database.SeedUsers()
	database.SeedSensors()

	// Tampilkan data Abdul Wahab dan Guest untuk verifikasi
	var wahab model.User
	database.DB.Where("nik = ?", "52.002.223").First(&wahab)
	fmt.Printf("[Verifikasi 1] Abdul Wahab: NIK=%s, Role=%s\n", wahab.GetNIK(), wahab.Role)

	var guest model.User
	database.DB.Where("role = ?", "Guest").First(&guest)
	nikVal := guest.NIK
	fmt.Printf("[Verifikasi 2] Guest: NIK=%s, PasswordHash='%s', Role=%s\n", nikVal, guest.PasswordHash, guest.Role)
}
