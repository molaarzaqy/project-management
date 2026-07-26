# Project Management API

Backend API untuk aplikasi manajemen proyek yang dibangun menggunakan bahasa **Go (Golang)** dengan framework **Fiber** dan menerapkan arsitektur **Service-Repository Pattern**.

---

## Architecture & Folder Structure

Proyek ini menggunakan **Service-Repository Pattern** untuk memisahkan *concern* (tanggung jawab) antar lapisan kode, sehingga kode lebih terstruktur, mudah diuji (*testable*), dan mudah di-maintain.

```text
project-management/
│
├── config/         # Konfigurasi aplikasi (Environment, DB connection)
├── controllers/    # Lapisan HTTP: Menerima request dan mengirim response
├── database/       # Manajemen database
│   ├── migrations/ # File migrasi SQL (schema up & down)
│   └── seed/       # Seeder data awal (misal: super admin)
├── docs/           # File dokumentasi Swagger (Generated - diabaikan oleh Git)
├── middleware/     # Middleware (JWT Auth, Error handling)
├── models/         # Definisi struct database (Entity) & DTO
│   └── types/      # Custom types untuk database (misal: UUIDArray untuk uuid[])
├── repositories/   # Lapisan Database (GORM query)
├── routes/         # Definisi dan pengelompokan Endpoint URL
├── services/       # Lapisan Logika Bisnis
├── utils/          # Fungsi pembantu (Token JWT, Hashing, Response)
├── .env            # Environment variables (lokal)
├── .gitignore      # Daftar file/folder yang diabaikan Git
├── go.mod          # Go module dependencies
├── go.sum          # Go checksum dependencies
└── main.go         # Entry point aplikasi

```

## API Documentation (Swagger)

Karena folder `docs/` diabaikan oleh Git (`.gitignore`), Anda perlu men-*generate* dokumentasi Swagger secara mandiri di komputer lokal Anda dengan langkah-langkah berikut:

1. **Install Tool Swag** (jika belum memilikinya):
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
2. **Generate Swagger Docs**:
    ```bash
   swag init
3. **Jalankan Aplikasi**:
    ```bash
   go run main.go
4. **Akses Swagger UI**:
    ```bash
   http://localhost:3030/swagger/index.html
