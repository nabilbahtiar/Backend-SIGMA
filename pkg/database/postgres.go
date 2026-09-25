package database

import (
	"fmt"
	"log"

	"server-room-auth/internal/config"
	"server-room-auth/internal/model"
	"server-room-auth/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.AppConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	fmt.Println("Koneksi database berhasil!")
}

func SeedSensors() {
	var count int64
	DB.Model(&model.Sensor{}).Count(&count)
	if count > 0 {
		return // Jika sudah ada data, jangan di-seed lagi
	}

	// Buat 6 Sensor
	sensors := []model.Sensor{
		{ID: uuid.New(), SensorCode: "SHT-01", SensorName: "Suhu Ruangan", SensorType: "TEMPERATURE", Unit: "°C", Location: "Ruang Server"},
		{ID: uuid.New(), SensorCode: "SHT-02", SensorName: "Kelembapan (RH)", SensorType: "HUMIDITY", Unit: "%", Location: "Ruang Server"},
		{ID: uuid.New(), SensorCode: "SMK-01", SensorName: "Proteksi Asap", SensorType: "SMOKE", Unit: "PPM", Location: "Plafon Tengah"},
		{ID: uuid.New(), SensorCode: "MAG-01", SensorName: "Keamanan Pintu", SensorType: "DOOR", Unit: "Boolean", Location: "Pintu Masuk"},
		{ID: uuid.New(), SensorCode: "VIB-01", SensorName: "Sensor Getaran", SensorType: "VIBRATION", Unit: "g", Location: "Rangka Rak"},
		{ID: uuid.New(), SensorCode: "WLK-01", SensorName: "Kebocoran Air", SensorType: "WATER_LEAK", Unit: "Boolean", Location: "Bawah AC"},
	}

	for i := range sensors {
		DB.Create(&sensors[i])
	}

	// Konfigurasi Threshold (Ambang Batas)
	configs := []model.SensorConfiguration{
		{ID: uuid.New(), SensorID: sensors[0].ID, MinThreshold: 18, MaxThreshold: 24, SeverityLevel: "WARNING"}, // Suhu Warning
		{ID: uuid.New(), SensorID: sensors[0].ID, MinThreshold: 0, MaxThreshold: 30, SeverityLevel: "CRITICAL"}, // Suhu Critical > 30 (Sesuai PDF)
		{ID: uuid.New(), SensorID: sensors[1].ID, MinThreshold: 40, MaxThreshold: 60, SeverityLevel: "WARNING"}, // Kelembapan
		{ID: uuid.New(), SensorID: sensors[2].ID, MinThreshold: 0, MaxThreshold: 10, SeverityLevel: "EMERGENCY"},// Asap (Emergency jika terdeteksi)
		// Untuk sensor boolean (pintu/air), kita pakai konvensi: misal 0 = Normal, 1 = Bahaya
		// Ini bisa disesuaikan nanti dengan logika backend
		{ID: uuid.New(), SensorID: sensors[3].ID, MinThreshold: 0, MaxThreshold: 0, SeverityLevel: "WARNING"},   // Pintu Terbuka
		{ID: uuid.New(), SensorID: sensors[4].ID, MinThreshold: 0, MaxThreshold: 0.05, SeverityLevel: "CRITICAL"},// Getaran Rak
		{ID: uuid.New(), SensorID: sensors[5].ID, MinThreshold: 0, MaxThreshold: 0, SeverityLevel: "CRITICAL"},  // Kebocoran Air
	}

	for i := range configs {
		DB.Create(&configs[i])
	}

	fmt.Println("Seeder: 6 Sensor IoT beserta Konfigurasinya berhasil disuntikkan!")
}

func SeedUsers() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return // Jika sudah ada data, jangan di-seed lagi
	}

	passDefault, _ := utils.HashPassword("rsch123")
	passGuest, _ := utils.HashPassword("guest123")

	users := []model.User{
		{ID: uuid.New(), NIK: "0421.00005", Nama: "ARIFIN EFENDI", Jabatan: "Security", Unit: "Security", TipePegawai: "internal", NoHP: "081233109975", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0621.00016", Nama: "RAHMAT RIYANTO", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "082139599950", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0621.00015", Nama: "MOCHAMAD RONI", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "085311048313", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0220.02202", Nama: "SUHARTONO", Jabatan: "Koordinator-Security", Unit: "Security", TipePegawai: "external", NoHP: "085784536927", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0222.00055", Nama: "MOCH. JAILANI", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "081238904343", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0222.00052", Nama: "HADI WIJAYA", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "082336726008", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0422.00067", Nama: "BELLA KURNIA CRISTA", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "085943441801", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0224.00096", Nama: "RISKI ADI PUTRA", Jabatan: "Security", Unit: "Security", TipePegawai: "internal", NoHP: "081252532093", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "1024.00109", Nama: "GESTI HOLILA", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "089530367474", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0125.00113", Nama: "Edo Candra Putra", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "0", Role: "Security", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0518.02163", Nama: "Rani Ekasari Pratiwi, Amd.", Jabatan: "IT Support", Unit: "IT", TipePegawai: "internal", NoHP: "082234514825", Role: "IT Support", PasswordHash: passDefault, IDTelegram: "@ekasari_pratiwi"},
		{ID: uuid.New(), NIK: "0715.02124", Nama: "Agung Sunaryo, S.Kom", Jabatan: "Penanggung Jawab Informasi dan Teknologi", Unit: "IT", TipePegawai: "internal", NoHP: "08990523963", Role: "Super Admin IT", PasswordHash: passDefault, IDTelegram: "@Wicwiky"},
		{ID: uuid.New(), NIK: "0523.02239", Nama: "Haris Arifin, S.Kom", Jabatan: "IT Support", Unit: "IT", TipePegawai: "internal", NoHP: "082338833248", Role: "IT Support", PasswordHash: passDefault, IDTelegram: "@harisrscitrahusada"},
		{ID: uuid.New(), NIK: "0624.02246", Nama: "Didit Purwanto", Jabatan: "Umum RT", Unit: "Umum RT", TipePegawai: "internal", NoHP: "082142846778", Role: "Facility/Engineering", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0416.02134", Nama: "Ali Ridho Arifi", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "082132222123", Role: "Facility/Engineering", PasswordHash: passDefault, IDTelegram: "@Ridhoarifi"},
		{ID: uuid.New(), NIK: "0220.02199", Nama: "M. Imron", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "08971543954", Role: "Facility/Engineering", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0120.02213", Nama: "Angga Prahanian Syah", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "082131539520", Role: "Facility/Engineering", PasswordHash: passDefault, IDTelegram: "@Anggaprahaniansyah"},
		{ID: uuid.New(), NIK: "52.002.223", Nama: "Abdul Wahab", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "085233151284", Role: "Facility/Engineering", PasswordHash: passDefault, IDTelegram: "@Bangwahab"},
		{ID: uuid.New(), NIK: "0816.02139", Nama: "Ageng Supriadi", Jabatan: "Ka. Unit-Umum RT", Unit: "Umum RT", TipePegawai: "internal", NoHP: "085233796252", Role: "Facility/Engineering", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0325.02251", Nama: "Dimas Adi Firmansyah", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "081276805711", Role: "Facility/Engineering", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0915.01133", Nama: "dr. Fatkhur Ruli Malik Qilsi", Jabatan: "Direktur", Unit: "Direksi", TipePegawai: "internal", NoHP: "081326992108", Role: "Manajemen", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0309.02117", Nama: "Andre Kartawidjaja, B.Sc", Jabatan: "Ka. Umum dan Keuangan", Unit: "Direksi", TipePegawai: "internal", NoHP: "081230351153", Role: "Manajemen", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "0317.01158", Nama: "dr. Dhea Anyssa Rachmati", Jabatan: "Ka. Bidang Yanmed", Unit: "Pelayanan Medik", TipePegawai: "internal", NoHP: "082244996959", Role: "Manajemen", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "102.402.322", Nama: "dr. Andritta Febriana, Sp.MK", Jabatan: "Ka. Bidang Jangmed", Unit: "Penunjang Medik", TipePegawai: "internal", NoHP: "081332019999", Role: "Manajemen", PasswordHash: passDefault, IDTelegram: ""},
		{ID: uuid.New(), NIK: "guest", Nama: "Guest", Jabatan: "-", Unit: "-", TipePegawai: "-", NoHP: "-", Role: "Guest", PasswordHash: passGuest, IDTelegram: ""},
	}

	for _, user := range users {
		DB.Create(&user)
	}

	fmt.Println("Seeder: 25 akun berhasil disuntikkan ke database (24 staf + 1 akun Guest)!")
}
