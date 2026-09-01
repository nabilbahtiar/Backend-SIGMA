# 🌡️ Server Room IoT Monitoring - Backend Service

Backend Service untuk Sistem Pemantauan & Keamanan Ruang Server On-Premise (RS Citra Husada Jember). Dibangun menggunakan **Golang** dengan framework **Gin** dan **GORM**.

## ✨ Fitur Utama (Fase 1: Authentication)
- **Autentikasi JWT**: Menggunakan token JWT (`golang-jwt/jwt/v5`) yang berlaku 24 jam.
- **Enkripsi Password**: Keamanan password menggunakan algoritma `bcrypt`.
- **Role-Based Access Control (RBAC)**: Middleware untuk membatasi akses endpoint berdasarkan *role* pengguna (Super Admin, IT Infrastructure, IT Support, dll).
- **Auto Seeder**: Tabel dan data pengguna awal (termasuk *superadmin*) akan dibuat secara otomatis saat program dijalankan.
- **CORS Enabled**: Sudah dikonfigurasi agar bisa diakses oleh *Frontend* dari port atau origin manapun (untuk *development*).

---

## 🛠️ Persyaratan Sistem (Prerequisites)
Pastikan Anda sudah menginstal:
- [Go (Golang)](https://go.dev/dl/) versi 1.20 atau lebih baru.
- [PostgreSQL](https://www.postgresql.org/download/).

---

## 🚀 Cara Menjalankan Proyek (Setup & Run)

**1. Clone Repository**
```bash
git clone https://github.com/<username-anda>/<nama-repo>.git
cd <nama-repo>
```

**2. Siapkan Database PostgreSQL**
- Buka pgAdmin atau psql.
- Buat sebuah database kosong baru. Contoh: `iot_server_room`.
- *Catatan: Anda tidak perlu membuat tabel secara manual, sistem akan membuatnya untuk Anda.*

**3. Sesuaikan Konfigurasi Koneksi DB**
- Buka file `main.go`.
- Cari baris `dsn := "host=localhost user=postgres password=123 dbname=iot_server_room port=5432 sslmode=disable TimeZone=Asia/Jakarta"` (berada di dalam fungsi `InitDB()`).
- Ubah `user` dan `password` sesuai dengan akun PostgreSQL di komputer Anda.

**4. Install Dependencies & Jalankan Server**
```bash
go mod tidy
go run main.go
```
*Server akan berjalan di `http://localhost:8080` dan data seeder otomatis dimasukkan ke database.*

---

## 📚 Dokumentasi API Terkini

### 1. Login (Mendapatkan Token)
- **Method:** `POST`
- **Endpoint:** `/api/login`
- **Body (JSON):**
  ```json
  {
      "username": "superadmin",
      "password": "password123"
  }
  ```
- **Response Sukses:** Mengembalikan Token JWT dan data *Role*.

### 2. Cek Status Dashboard (Semua Role)
- **Method:** `GET`
- **Endpoint:** `/api/secure/dashboard/status`
- **Header:** `Authorization: Bearer <token_jwt>`

### 3. Konfigurasi Sensor (Khusus Admin)
- **Method:** `POST`
- **Endpoint:** `/api/secure/sensor/config`
- **Header:** `Authorization: Bearer <token_jwt>`
- *(Akan mengembalikan error 403 Forbidden jika role bukan Super Admin / IT Infrastructure Admin).*

---

*Dibangun berdasarkan Dokumen Spesifikasi Teknis & Blueprint IoT Server Room RS Citra Husada Jember.*
