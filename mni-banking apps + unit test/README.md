Dokumentasi Aplikasi Mini Banking
Aplikasi ini adalah layanan RESTful API sederhana yang dibangun menggunakan Go (Golang) untuk simulasi perbankan. Aplikasi ini menggunakan framework Gin untuk routing HTTP dan berinteraksi dengan database PostgreSQL.

Daftar Isi
Struktur Proyek
Arsitektur
Konfigurasi
Model Data
Lapisan-lapisan Aplikasi
Controller
Usecase
Repository
API Endpoints
Setup dan Menjalankan Aplikasi
Menjalankan Unit Test
Struktur Proyek
plaintext
 Show full code block 
/
├── config/                 # Konfigurasi aplikasi (database, JWT, dll.)
├── controller/             # Handler HTTP (Gin)
├── delivery/               # Entry point aplikasi (server HTTP)
├── middleware/             # Middleware untuk request (misal: otentikasi)
├── model/                  # Struct untuk model data (User, Transaction)
├── repository/             # Logika akses database
├── usecase/                # Logika bisnis aplikasi
├── utils/                  # Utilitas (misal: service JWT)
├── go.mod                  # Manajemen dependensi Go
├── go.sum
├── main.go                 # File utama untuk menjalankan server
└── README.md
Arsitektur
Aplikasi ini mengadopsi arsitektur berlapis (Layered Architecture) untuk memisahkan tanggung jawab dan meningkatkan modularitas.

Controller Layer: Bertanggung jawab untuk menerima permintaan HTTP, memvalidasi input, dan memanggil use case yang sesuai. Lapisan ini berinteraksi langsung dengan klien (misalnya, browser atau aplikasi mobile).
Usecase Layer (Business Logic): Berisi logika bisnis inti dari aplikasi. Lapisan ini tidak bergantung pada detail implementasi seperti database atau framework web. Misalnya, proses registrasi pengguna yang melibatkan hashing password ada di sini.
Repository Layer (Data Access): Bertanggung jawab untuk berinteraksi dengan database. Lapisan ini mengabstraksi operasi CRUD (Create, Read, Update, Delete) ke database.
Alur permintaan (request flow) berjalan sebagai berikut: Request HTTP -> Middleware -> Controller -> Usecase -> Repository -> Database

Konfigurasi
Konfigurasi aplikasi diatur melalui environment variables dan dikelola oleh paket config. File .env (tidak ada dalam konteks) dapat digunakan untuk pengembangan lokal.

Variabel yang perlu diatur:

DB_HOST: Host database (default: localhost)
DB_PORT: Port database (default: 5432)
DB_NAME: Nama database
DB_USER: Username database
DB_PASSWORD: Password database (wajib)
API_PORT: Port untuk server API (default: 8080)
JWT_SIGNATURE_KEY: Kunci rahasia untuk menandatangani token JWT (wajib)
Model Data
Aplikasi ini memiliki dua model utama:

UserCredential: Merepresentasikan pengguna.
Id (uint32): ID unik pengguna.
Username (string): Nama pengguna untuk login.
Password (string): Password yang sudah di-hash.
Role (string): Peran pengguna (misal: 'admin', 'user').
Balance (float64): Saldo akun pengguna.

Transaction: Merepresentasikan transaksi antar pengguna.
ID (uint): ID unik transaksi.
FromUserID (uint32): ID pengguna pengirim.
ToUserID (uint32): ID pengguna penerima.
Amount (float64): Jumlah yang ditransfer.
Type (string): Jenis transaksi (misal: 'transfer').
Status (string): Status transaksi (misal: 'completed').

Lapisan-lapisan Aplikasi
Controller
Terletak di direktori controller/. Menggunakan framework Gin untuk menangani rute HTTP.

UserController:
createUser: Menangani registrasi pengguna baru (POST /api/v1/users).
getAllUser: Mengambil daftar semua pengguna (GET /api/v1/users).
getUserById: Mengambil detail pengguna berdasarkan ID (GET /api/v1/users/:id).
Usecase
Terletak di direktori usecase/. Berisi logika bisnis.

userUseCase:
RegisterNewUser: Memproses pendaftaran pengguna baru, termasuk hashing password sebelum menyimpannya.
FindAllUser: Mengambil semua data pengguna dari repository.
FindUserById: Mencari pengguna berdasarkan ID.
FindUserByUsernamePassword: Memverifikasi kredensial pengguna saat login. Ini akan mengambil data pengguna berdasarkan username dari repository, lalu membandingkan hash password yang tersimpan dengan password yang diberikan menggunakan bcrypt.
Repository
Terletak di direktori repository/. Menangani semua interaksi dengan database.

userRepository:
Create: Menyimpan pengguna baru ke tabel mst_user.
List: Mengambil semua pengguna dari database.
Get: Mengambil satu pengguna berdasarkan ID.
GetByUsername: Mengambil satu pengguna berdasarkan username. Fungsi ini penting untuk proses login dan memeriksa apakah username sudah ada.
transactionRepository:

Create: Membuat entri transaksi baru di tabel mst_transaction.
List: Mengambil semua transaksi.
Get: Mengambil transaksi berdasarkan ID.
GetByUserId: Mengambil semua transaksi yang melibatkan ID pengguna tertentu (baik sebagai pengirim maupun penerima).
Update: Memperbarui data transaksi.
Delete: Menghapus transaksi dari database.
API Endpoints
Semua endpoint berada di bawah prefix /api/v1.

Method	Endpoint	Deskripsi	Otentikasi
POST	/users	Mendaftarkan pengguna baru.	Tidak perlu
GET	/users	Mendapatkan daftar semua pengguna.	Token JWT (role 'admin')
GET	/users/:id	Mendapatkan detail pengguna berdasarkan ID.	Token JWT (role 'admin')
Catatan: Endpoint untuk login dan transaksi perlu ditambahkan untuk fungsionalitas penuh.

Setup dan Menjalankan Aplikasi
Database Setup: Pastikan Anda memiliki database PostgreSQL yang berjalan. Buat database dan tabel menggunakan skrip SQL dari README.md.

sql
 Show full code block 
CREATE DATABASE your_db_name;

\c your_db_name;

CREATE TABLE mst_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    role VARCHAR(100),
    balance NUMERIC
);

CREATE TABLE mst_transaction (
    id SERIAL PRIMARY KEY,
    from_user_id INT,
    to_user_id INT,
    amount NUMERIC,
    type VARCHAR(50),
    status VARCHAR(50)
);
Konfigurasi Environment: Buat file .env di root proyek dan isi dengan konfigurasi database Anda.

plaintext
 Show full code block 
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_db_name
DB_USER=your_username
DB_PASSWORD=your_password
API_PORT=8080
JWT_SIGNATURE_KEY=your_super_secret_key
Install Dependensi: Jalankan perintah berikut di terminal.

bash
go mod tidy
Jalankan Aplikasi:

bash
go run main.go
Server akan berjalan di port yang ditentukan (default: 8080).

Menjalankan Unit Test
Aplikasi ini dilengkapi dengan unit test untuk setiap lapisan menggunakan testify dan sqlmock.

Untuk menjalankan semua tes:

bash
go test ./...

Untuk menjalankan tes pada paket tertentu (contoh: repository):

bash
go test ./repository
Dokumentasi ini memberikan gambaran tingkat tinggi tentang aplikasi. Untuk detail implementasi, Anda dapat merujuk langsung ke kode sumber di setiap paket.
