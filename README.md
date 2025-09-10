# Project GoLang

Repositori ini adalah aplikasi backend menggunakan **Go**. Berikut adalah struktur dan panduan setup dasar.

---

## 📁 Struktur Folder
```bash
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

```


---

## ⚙️ Prerequisites

Sebelum memulai, pastikan sudah terinstall:

- [Go](https://go.dev/) >= 1.23.3
- [PostgreSQL](https://www.postgresql.org/) (atau database lain sesuai config)
- Git

---

## 📦 Dependencies

Project ini menggunakan beberapa library Go yang dikelola lewat `go.mod`:

- [Fiber](https://gofiber.io/) → web framework
- [Gorm](https://gorm.io/) → ORM untuk database
- [Validator](https://github.com/go-playground/validator) → validasi struct
- [UUID](https://github.com/google/uuid) → generate UUID
- [JWT](https://github.com/golang-jwt/jwt) → autentikasi token

---

## 🗄️ Database ERD

Skema database project ini dapat dilihat melalui gambar atau link berikut:

### Link ERD
[Schema ERD](https://dbdiagram.io/d/68c0d94561a46d388e4b20d0)

### Preview ERD
![ERD Database](<img width="1024" height="974" alt="schema ERD" src="https://github.com/user-attachments/assets/871b2a44-fbb7-4a28-9b29-2554693fc764" />)

---

## 🚀 Panduan Pengaturan Proyek

1.  **Klon Repositori**

    Buka terminal dan jalankan perintah berikut:
    ```bash
    git clone [https://github.com/username/repo.git](https://github.com/username/repo.git)
    cd repo
    ```

2.  **Konfigurasi Variabel Lingkungan**

    Salin file `.env.example` menjadi `.env` dan sesuaikan nilainya:
    ```bash
    cp .env.example .env
    ```

3.  **Instalasi Dependensi**

    Jalankan perintah ini untuk mengunduh semua modul yang diperlukan:
    ```bash
    go mod tidy
    ```

4.  **Menjalankan Aplikasi**

    Jalankan aplikasi dengan perintah:
    ```bash
    go run main.go
    ```
    Aplikasi akan berjalan pada port yang ditentukan di file `.env` (port default `3001`).
```

