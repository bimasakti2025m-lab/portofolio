📄 Dokumentasi Proyek: Go JWT Authentication dengan Clean Architecture
Aplikasi ini merupakan contoh implementasi autentikasi menggunakan JSON Web Token (JWT) pada bahasa pemrograman Go. Proyek ini dibangun dengan menerapkan prinsip-prinsip Clean Architecture untuk memisahkan lapisan-lapisan aplikasi, sehingga lebih mudah untuk dikelola, diuji, dan dikembangkan.

✨ Fitur Utama
Autentikasi Pengguna: Proses login untuk mendapatkan token JWT.

Middleware Autentikasi: Melindungi endpoint tertentu yang hanya bisa diakses dengan token JWT yang valid.

Clean Architecture: Struktur proyek yang terorganisir, memisahkan lapisan Delivery (Controller), Use Case (Service), dan Repository.

Unit Testing: Dilengkapi dengan unit test (go test ./... -v) untuk memastikan setiap fungsi berjalan sesuai harapan.

Manajemen Konfigurasi: Menggunakan file .env untuk mengelola variabel lingkungan seperti JWT_SECRET dan port server.

🏛️ Struktur Arsitektur (Clean Architecture)
Proyek ini mengadopsi prinsip Clean Architecture untuk memisahkan tanggung jawab (Separation of Concerns).

Direktori	Nama Lapisan	Tanggung Jawab Utama
**/controller	Delivery/Frameworks	Menangani permintaan HTTP dan memberikan respons. Meneruskan permintaan ke lapisan service.
**/service	Use Case/Application	Berisi semua logika bisnis aplikasi. Bergantung pada interface dari repository.
**/repository	Adapters	Berinteraksi dengan sumber data (Database, Mock Data, dll.).
**/model	Entities	Berisi struktur data (struct) yang digunakan di seluruh aplikasi (DTOs, Domain Models).
**/middleware	Frameworks	Berisi fungsi-fungsi seperti validasi token JWT sebelum mencapai controller.
main.go	Main	Titik masuk utama aplikasi, menginisialisasi semua komponen dan menjalankan server.

Export to Sheets

🚀 Cara Menjalankan Aplikasi
Prasyarat
Go versi 1.18 atau lebih baru.

1. Clone Repositori
git clone <URL_REPOSITORI_ANDA>
cd <NAMA_DIREKTORI_PROYEK>

2. Konfigurasi Environment
Buat file bernama .env di direktori utama proyek. Sangat disarankan untuk menyertakan file .env.example di repositori.

Code snippet

# Port untuk server HTTP
PORT=8080

# Kunci rahasia untuk menandatangani token JWT (Ganti dengan kunci yang kuat!)
JWT_SECRET=kunci_rahasia_anda_yang_sangat_aman

3. Unduh Dependensi
Jalankan perintah berikut untuk mengunduh semua modul yang dibutuhkan:
go mod tidy

4. Jalankan Aplikasi
Gunakan perintah di bawah ini untuk menjalankan server HTTP:

go run main.go
Server akan berjalan di http://localhost:8080 (atau port yang Anda tentukan di .env).

🧪 Menjalankan Unit Test
Untuk memastikan semua fungsi (terutama di lapisan service dan repository) berjalan dengan benar, jalankan unit test dengan perintah berikut:
go test ./... -v
🌐 Endpoints API
Berikut adalah daftar endpoint API yang tersedia di aplikasi ini:

1. Login Pengguna
Endpoint: POST /login

Deskripsi: Mengautentikasi pengguna dan mengembalikan token JWT.

Parameter	Tipe	Deskripsi
username	string	Nama pengguna untuk login.
password	string	Kata sandi pengguna.

Export to Sheets

Request Body (JSON):

JSON

{
  "username": "user1",
  "password": "password123"
}
Response Sukses (200 OK):

JSON

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

2. Mengakses Rute Terproteksi (Protected Route)
Endpoint: GET /profile

Deskripsi: Contoh endpoint yang dilindungi oleh middleware JWT. Membutuhkan Authorization header yang valid.

Headers:

Authorization: Bearer <TOKEN_JWT_ANDA>
Response Sukses (200 OK):

JSON

{
  "message": "Welcome user1"
}
Response Gagal (401 Unauthorized):

JSON

{
  "error": "Invalid token"
}
💡 Saran Tambahan
URL Repositori: Pastikan untuk mengganti <URL_REPOSITORI_ANDA> dan <NAMA_DIREKTORI_PROYEK> dalam dokumentasi clone repositori di atas.

# QUERY MEMBUAT DATABASE DAN TABEL
1. CREATE DATABASE 
    CREATE DATABASE users_db;

2. CREATE TABLE
    -- Disarankan untuk membuat database secara terpisah terlebih dahulu
-- CREATE DATABASE users_db;

-- --- Hubungkan ke database users_db sebelum menjalankan query di bawah ---

-- Membuat tabel 'users' untuk menyimpan data pengguna
CREATE TABLE users (
    -- ID unik untuk setiap pengguna, akan bertambah otomatis
    id SERIAL PRIMARY KEY,

    -- Username pengguna, harus unik dan tidak boleh kosong
    username VARCHAR(50) UNIQUE NOT NULL,

    -- Password pengguna. Di aplikasi nyata, kolom ini HARUS menyimpan password yang sudah di-hash!
    password VARCHAR(255) NOT NULL,

    -- Role pengguna
    role VARCHAR(50) DEFAULT 'user',
);

-- KOMENTAR PENTING:
-- 1. Keamanan Password: Jangan pernah menyimpan password dalam bentuk teks biasa (plain text).
--    Sebelum menyimpan ke database, selalu hash password menggunakan algoritma yang kuat seperti bcrypt atau Argon2.
--    Panjang VARCHAR(255) sudah cukup untuk menampung hasil hash dari kebanyakan algoritma.
--
-- 2. Username Unik: Constraint `UNIQUE` pada kolom `username` memastikan tidak ada dua pengguna
--    dengan username yang sama, yang penting untuk proses login.


-- (Opsional) Query untuk memasukkan data contoh
-- Ganti 'hash_dari_password123' dengan hasil hash password yang sebenarnya.
INSERT INTO users (username, password) VALUES
('user1', '$2a$10$your_bcrypt_hash_for_password123_here');

