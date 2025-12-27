// ...existing code...
# E-commerce Sederhana — Dokumentasi Lengkap

Ringkasan singkat  
Aplikasi backend sederhana menggunakan Go dengan arsitektur terpisah (Clean Architecture). Fitur utama: autentikasi JWT, manajemen produk, keranjang, pesanan, dan integrasi notifikasi Midtrans.

Repository & simbol utama
- File entry point: [`main.go`](main.go) (fungsi [`main`](main.go))
- Server builder: [`main.NewServer`](server.go) — inisialisasi dependensi dan route
- Konfigurasi: [`config.NewConfig`](config/config.go) — baca env dan konfigurasinya
- Middleware JWT: [`middleware.AuthMiddleware`](middleware/auth_middleware.go)
- JWT service: [`service.JWTservice`](utils/service/jwt_service.go) (`CreateToken`, `VerifyToken`)
- Midtrans service & handler: [`midtrans.MidtransService`](utils/service/midtrans/midtrans_service.go) dan [`midtrans.MidtransHandler`](utils/service/midtrans/midtrans_handler.go)
- Usecase contoh: [`usecase.UserUsecase`](usecase/user_usecase.go), [`usecase.OrderUsecase`](usecase/order_usecase.go)
- Repository contoh: [`repository.UserRepository`](repository/user_repository.go), [`repository.OrderRepository`](repository/order_repository.go)
- Model JWT claims: [`modelutils.JwtPayloadClaims`](utils/model_utils/jwt_model.go)

Struktur proyek (intuitif)
- [controller/](controller) — HTTP handlers (mis. [`controller.UserController`](controller/user_controller.go))
- [usecase/](usecase) — logika bisnis / use cases
- [repository/](repository) — akses DB
- [model/](model) — entitas domain
- [middleware/](middleware) — auth middleware
- [utils/service/](utils/service) — jwt & midtrans integrations
- [config/](config) — pembacaan .env & konfigurasi

Persiapan (local)
1. Salin file `.env` (contoh) atau buat sendiri:
   - DB_HOST
   - DB_PORT
   - DB_USERNAME
   - DB_PASSWORD
   - DB_DATABASE
   - DB_DRIVER (ex: postgres)
   - API_PORT (contoh: 8080)
   - TOKEN_APPLICATION_NAME
   - TOKEN_JWT_SIGNATURE_KEY
   - TOKEN_ACCESS_TOKEN_LIFETIME (contoh: 1h)
   - MIDTRANS_SERVER_KEY
   - MIDTRANS_ENV (sandbox|production)

   Konfigurasi dibaca oleh [`config.readConfig`](config/config.go).

2. Install deps:
```sh
go mod tidy

// ...existing code...
# E-commerce Sederhana — Dokumentasi Lengkap

Ringkasan singkat  
Aplikasi backend sederhana menggunakan Go dengan arsitektur terpisah (Clean Architecture). Fitur utama: autentikasi JWT, manajemen produk, keranjang, pesanan, dan integrasi notifikasi Midtrans.

Repository & simbol utama
- File entry point: [`main.go`](main.go) (fungsi [`main`](main.go))
- Server builder: [`main.NewServer`](server.go) — inisialisasi dependensi dan route
- Konfigurasi: [`config.NewConfig`](config/config.go) — baca env dan konfigurasinya
- Middleware JWT: [`middleware.AuthMiddleware`](middleware/auth_middleware.go)
- JWT service: [`service.JWTservice`](utils/service/jwt_service.go) (`CreateToken`, `VerifyToken`)
- Midtrans service & handler: [`midtrans.MidtransService`](utils/service/midtrans/midtrans_service.go) dan [`midtrans.MidtransHandler`](utils/service/midtrans/midtrans_handler.go)
- Usecase contoh: [`usecase.UserUsecase`](usecase/user_usecase.go), [`usecase.OrderUsecase`](usecase/order_usecase.go)
- Repository contoh: [`repository.UserRepository`](repository/user_repository.go), [`repository.OrderRepository`](repository/order_repository.go)
- Model JWT claims: [`modelutils.JwtPayloadClaims`](utils/model_utils/jwt_model.go)

Struktur proyek (intuitif)
- [controller/](controller) — HTTP handlers (mis. [`controller.UserController`](controller/user_controller.go))
- [usecase/](usecase) — logika bisnis / use cases
- [repository/](repository) — akses DB
- [model/](model) — entitas domain
- [middleware/](middleware) — auth middleware
- [utils/service/](utils/service) — jwt & midtrans integrations
- [config/](config) — pembacaan .env & konfigurasi

Persiapan (local)
1. Salin file `.env` (contoh) atau buat sendiri:
   - DB_HOST
   - DB_PORT
   - DB_USERNAME
   - DB_PASSWORD
   - DB_DATABASE
   - DB_DRIVER (ex: postgres)
   - API_PORT (contoh: 8080)
   - TOKEN_APPLICATION_NAME
   - TOKEN_JWT_SIGNATURE_KEY
   - TOKEN_ACCESS_TOKEN_LIFETIME (contoh: 1h)
   - MIDTRANS_SERVER_KEY
   - MIDTRANS_ENV (sandbox|production)

   Konfigurasi dibaca oleh [`config.readConfig`](config/config.go).

2. Install deps:
```sh
go mod tidy

3. Jalankan:
```sh
go run .

Server menjalankan host ":" + cfg.API.Port — inisialisasi di main.NewServer.

Database — contoh schema
(diturunkan dari repository implementations)

● users
  ○ id SERIAL PRIMARY KEY
  ○ username VARCHAR UNIQUE NOT NULL
  ○ email VARCHAR
  ○ password VARCHAR NOT NULL
  ○ role VARCHAR DEFAULT 'user'
● products
  ○ id, name, description, price (float), stock (int)
● carts, cart_items
● orders, order_items

Contoh SQL ringkas :
  CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user'
  );

  CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL
)

  CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status status NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
)

  CREATE TABLE cart_items (
      id SERIAL PRIMARY KEY,
      cart_id INT NOT NULL,
      product_id INT NOT NULL,
      quantity INT NOT NULL,
      FOREIGN KEY (cart_id) REFERENCES carts(id),
      FOREIGN KEY (product_id) REFERENCES products(id)
  )

  CREATE TABLE orders(
      id SERIAL PRIMARY KEY,
      user_id INT NOT NULL,
      total DECIMAL(10, 2) NOT NULL,
      status_pesanan VARCHAR(50) NOT NULL,
      transaction_id_midtrans VARCHAR(255),
      FOREIGN KEY (user_id) REFERENCES users(id)
  )

  CREATE TABLE order_items(
      id SERIAL PRIMARY KEY,
      order_id INT NOT NULL,
      product_id INT NOT NULL,
      quantity INT NOT NULL,
      price_snapshot DECIMAL(10, 2),
      FOREIGN KEY (order_id) REFERENCES orders(id),
      FOREIGN KEY (product_id) REFERENCES products(id)
  )

Endpoint API (prefix /api/v1)

● Auth
  ○ POST /api/v1/register -> controller.AuthController.registerHandler
  ○ POST /api/v1/login -> controller.AuthController.loginHandler
● User (admin-only via middleware.AuthMiddleware.RequireToken)
  ○ POST /api/v1/users -> controller.UserController.createUserHandler
  ○ GET /api/v1/users -> controller.UserController.getAllUsersHandler
  ○ GET /api/v1/users/:username -> controller.UserController.getUserByUsernameHandler
● Product (admin-only)
  ○ CRUD: controller.ProductController
● Cart & Cart Items (user)
  ○ controller.CartController
  ○ controller.CartItemController
● Order & Order Items (user)
controller.OrderController
controller.OrderItemController
  ○ Integrasi Midtrans:
    ■ Creation flow: usecase.OrderUsecase.CreateOrder memanggil midtrans.MidtransService.CreateTransaction
    ■ Notifikasi webhook: POST /api/v1/midtrans/notification -> midtrans.MidtransHandler.HandleNotification

### API Documentation

Dokumentasi ini menjelaskan setiap endpoint yang tersedia di dalam API. Semua endpoint berada di bawah prefix `/api/v1`.

#### Authentication (`/auth`)

**1. Register User**
- **Endpoint:** `POST /api/v1/register`
- **Deskripsi:** Mendaftarkan pengguna baru.
- **Request Body:**
  ```json
  {
    "username": "newuser",
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Success Response (201 Created):**
  ```json
  {
    "message": "User created successfully",
    "data": {
      "id": 1,
      "username": "newuser",
      "email": "user@example.com",
      "role": "user"
    }
  }
  ```

**2. Login User**
- **Endpoint:** `POST /api/v1/login`
- **Deskripsi:** Mengautentikasi pengguna dan mengembalikan token JWT.
- **Request Body:**
  ```json
  {
    "username": "newuser",
    "password": "password123"
  }
  ```
- **Success Response (200 OK):**
  ```json
  {
    "token": "ey..."
  }
  ```

---

#### Products (`/products`)
- **Otorisasi:** Memerlukan token `Bearer` dengan role `admin`.

**1. Create Product**
- **Endpoint:** `POST /api/v1/products`
- **Deskripsi:** Menambahkan produk baru.
- **Request Body:**
  ```json
  {
    "name": "Laptop Pro",
    "description": "Laptop canggih untuk profesional",
    "price": 15000000.00,
    "stock": 50
  }
  ```
- **Success Response (201 Created):** Mengembalikan data produk yang baru dibuat.

**2. Get All Products**
- **Endpoint:** `GET /api/v1/products`
- **Deskripsi:** Mendapatkan daftar semua produk.
- **Success Response (200 OK):**
  ```json
  {
    "data": [
      {
        "id": 1,
        "name": "Laptop Pro",
        "description": "Laptop canggih untuk profesional",
        "price": 15000000.00,
        "stock": 50
      }
    ]
  }
  ```

---

#### Cart (`/cart`)
- **Otorisasi:** Memerlukan token `Bearer` dengan role `user`.

**1. Add Item to Cart**
- **Endpoint:** `POST /api/v1/cart/items`
- **Deskripsi:** Menambahkan item ke keranjang belanja pengguna yang sedang aktif.
- **Request Body:**
  ```json
  {
    "product_id": 1,
    "quantity": 2
  }
  ```
- **Success Response (201 Created):** Mengembalikan detail item yang ditambahkan.

---

#### Orders (`/orders`)
- **Otorisasi:** Memerlukan token `Bearer` dengan role `user`.

**1. Create Order**
- **Endpoint:** `POST /api/v1/orders`
- **Deskripsi:** Membuat pesanan dari item yang ada di keranjang dan mengembalikan URL pembayaran Midtrans.
- **Success Response (201 Created):**
  ```json
  {
    "message": "Order created successfully, please proceed to payment",
    "data": {
      "payment_url": "https://app.sandbox.midtrans.com/snap/v1/transactions/..."
    }
  }
  ```


JWT & Middleware

● Token dibuat oleh service.JWTservice.CreateToken dan diverifikasi oleh service.JWTservice.VerifyToken.
● Claims custom pada payload menggunakan struct modelutils.JwtPayloadClaims.
● Middleware middleware.AuthMiddleware.RequireToken membaca header Authorization: "Bearer <token>" lalu memvalidasi role (multiple roles didukung).

Perubahan penting pada kode

● Untuk melihat inisialisasi server dan dependensi, buka server.go.
● Untuk melihat implementasi repository user: repository/user_repository.go.
● Untuk melihat alur pembuatan order + midtrans: usecase/order_usecase.go dan utils/service/midtrans/midtrans_service.go.

Kontribusi
● Ikuti pattern clean architecture: controller -> usecase -> repository.
● Tambahkan unit tests untuk tiap perubahan.
● Jangan commit file .env (terdaftar di .gitignore).

Kontak / sumber

● Baca kode utama untuk detail: controller/, usecase/, repository/, utils/service/, config/.

Disclaimer singkat
Dokumentasi ini ringkasan teknis untuk pengembangan cepat. Untuk produksi, perkuat validasi input, logging, error handling, dan keamanan (secret management, TLS, rate limiting).
