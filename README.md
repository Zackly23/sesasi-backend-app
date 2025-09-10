# Project GoLang

Project ini adalah aplikasi backend menggunakan **Go**. Berikut adalah struktur dan panduan setup dasar.

---

## 📁 Struktur Folder
```bash
# Example:
# Original tree command:

$ tree
.
├── dir1
│   ├── file11.ext
│   └── file12.ext
├── dir2
│   ├── file21.ext
│   ├── file22.ext
│   └── file23.ext
├── dir3
├── file_in_root.ext
└── README.md
```


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
