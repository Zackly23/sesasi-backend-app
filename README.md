# Project GoLang

Project ini adalah aplikasi backend menggunakan **Go**. Berikut adalah struktur dan panduan setup dasar.

---

## 📁 Struktur Folder
$ tree
.
├── config/ # Konfigurasi aplikasi (database, dll)
├── handlers/ # Handler untuk endpoint HTTP
├── middlewares/ # Middleware (auth, user)
├── models/ # Model database / struct
├── routes/ # Definisi route dan grouping endpoint
├── schemas/ # Schema validation / request & response structs
├── seeders/ # Script untuk seed data awal
├── utils/ # Utility functions / helpers
├── main.go # Entry point aplikasi
├── go.mod # Modul dependencies Go
├── go.sum # Checksum dependencies Go
├── .env.example # Contoh environment variables
├── .gitignore # File yang di-ignore oleh Git


---

## ⚙️ Prerequisites

- Go >= 1.21
- PostgreSQL / MySQL / database lain (sesuaikan dengan config)
- Git

---

## 🚀 Setup Project

1. **Clone repository**

```bash
git clone https://github.com/username/repo.git
cd repo
Copy .env.example menjadi .env dan isi sesuai environment kamu

bash
Copy code
cp .env.example .env
Install dependencies

bash
Copy code
go mod tidy
Jalankan aplikasi

bash
Copy code
go run main.go
Aplikasi akan berjalan di port yang sudah ditentukan di .env (default bisa 8080).

🗂 Folder Descriptions
config: berisi konfigurasi database, server, dan environment.

handlers: fungsi-fungsi yang menangani request HTTP.

middlewares: filter request, seperti autentikasi dan logging.

models: struct untuk database dan ORM mapping.

routes: define endpoint dan group route.

schemas: validation schema untuk request/response.

seeders: script untuk memasukkan data awal ke database.

utils: fungsi pembantu umum.

📦 Deploy / Production
Pastikan .env sudah terisi environment production.

Build aplikasi:

bash
Copy code
go build -o app
Jalankan binary:

bash
Copy code
./app
