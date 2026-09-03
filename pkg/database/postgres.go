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

	DB.AutoMigrate(&model.User{})
	fmt.Println("Koneksi database berhasil!")

	SeedUsers()
}

func SeedUsers() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return // Jika sudah ada data, jangan di-seed lagi
	}

	hashedPassword, _ := utils.HashPassword("rsch123")

	users := []model.User{
		{ID: uuid.New(), NIK: "0421.00005", Nama: "ARIFIN EFENDI", Jabatan: "Security", Unit: "Security", TipePegawai: "internal", NoHP: "081233109975", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0621.00016", Nama: "RAHMAT RIYANTO", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "82139599950", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0621.00015", Nama: "MOCHAMAD RONI", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "85311048313", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0220.02202", Nama: "SUHARTONO", Jabatan: "Koordinator-Security", Unit: "Security", TipePegawai: "external", NoHP: "85784536927", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0222.00055", Nama: "MOCH. JAILANI", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "81238904343", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0222.00052", Nama: "HADI WIJAYA", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "82336726008", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0422.00067", Nama: "BELLA KURNIA CRISTA", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "85943441801", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0224.00096", Nama: "RISKI ADI PUTRA", Jabatan: "Security", Unit: "Security", TipePegawai: "internal", NoHP: "81252532093", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "1024.00109", Nama: "GESTI HOLILA", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "89530367474", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0125.00113", Nama: "Edo Candra Putra", Jabatan: "Security", Unit: "Security", TipePegawai: "external", NoHP: "0", Role: "Petugas Jaga Keamanan", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0518.02163", Nama: "Rani Ekasari Pratiwi, Amd.", Jabatan: "IT", Unit: "IT", TipePegawai: "internal", NoHP: "082234514825", Role: "Petugas TIK", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0715.02124", Nama: "Agung Sunaryo, S.Kom", Jabatan: "Koordinator Informasi dan Teknologi", Unit: "IT", TipePegawai: "internal", NoHP: "08990523963", Role: "Superadmin", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0523.02239", Nama: "Haris Arifin, S.Kom", Jabatan: "IT", Unit: "IT", TipePegawai: "internal", NoHP: "082338833248", Role: "Petugas TIK", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0624.02246", Nama: "Didit Purwanto", Jabatan: "Umum RT", Unit: "Umum RT", TipePegawai: "internal", NoHP: "082142846778", Role: "Petugas Sapras", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0416.02134", Nama: "Ali Ridho Arifi", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "082132222123", Role: "Petugas Sapras", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0220.02199", Nama: "M. Imron", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "08971543954", Role: "Petugas Sapras", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0120.02213", Nama: "Angga Prahanian Syah", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "082131539520", Role: "Petugas Sapras", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0816.02139", Nama: "Ageng Supriadi", Jabatan: "Ka. Unit-Umum RT", Unit: "Umum RT", TipePegawai: "internal", NoHP: "085233796252", Role: "Koordinator Sapras", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0325.02251", Nama: "Dimas Adi Firmansyah", Jabatan: "Umum RT-IPSRS", Unit: "IPSRS", TipePegawai: "internal", NoHP: "081276805711", Role: "Petugas Sapras", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0915.01133", Nama: "dr. Fatkhur Ruli Malik Qilsi", Jabatan: "Direktur", Unit: "Direksi", TipePegawai: "internal", NoHP: "81326992108", Role: "Management", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0309.02117", Nama: "Andre Kartawidjaja, B.Sc", Jabatan: "Ka. Umum dan Keuangan", Unit: "Direksi", TipePegawai: "internal", NoHP: "81230351153", Role: "Management", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "0317.01158", Nama: "dr. Dhea Anyssa Rachmati", Jabatan: "Ka. Bidang Yanmed", Unit: "Pelayanan Medik", TipePegawai: "internal", NoHP: "82244996959", Role: "Management", PasswordHash: hashedPassword},
		{ID: uuid.New(), NIK: "102.402.322", Nama: "dr. Andritta Febriana, Sp. MK", Jabatan: "Ka. Bidang Jangmed", Unit: "Penunjang Medik", TipePegawai: "internal", NoHP: "081332019999", Role: "Management", PasswordHash: hashedPassword},
	}

	for _, user := range users {
		DB.Create(&user)
	}

	fmt.Println("Seeder: 23 akun baru berhasil disuntikkan ke database!")
}
