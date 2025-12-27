# URL Shortener Service (Golang)

Aplikasi REST API untuk memendekkan URL (URL Shortener) yang dibangun menggunakan Golang (Gin Framework) dan PostgreSQL. Aplikasi ini dilengkapi dengan otentikasi menggunakan JWT (JSON Web Token).

## Fitur

- **Otentikasi**: Registrasi dan Login pengguna (Admin/User).
- **URL Shortening**: Membuat kode pendek untuk URL panjang.
- **Lookup**: Mendapatkan URL asli berdasarkan kode pendek.
- **Unit Testing**: Dilengkapi dengan unit test untuk repository dan server.

## Persiapan Database

Buat database dan tabel menggunakan query SQL berikut. Pastikan PostgreSQL sudah terinstall.

```sql
-- Buat Database
CREATE DATABASE url_shortener_db;

-- Tabel User
CREATE TABLE mst_user
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role     VARCHAR(100) NOT NULL
);

-- Tabel URLs
CREATE TABLE mst_urls
(
    id SERIAL PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_code VARCHAR(50) UNIQUE NOT NULL,
    user_id INTEGER REFERENCES mst_user(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Data Awal (Seeding User)
INSERT INTO mst_user (username, password, role)
VALUES
    ('admin', 'admin', 'admin'),
    ('user', 'user', 'user');
```
