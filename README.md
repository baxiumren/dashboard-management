# 📊 FB Management Account Dashboard

<div align="center">

**Dashboard manajemen akun Facebook Ads — dibangun dengan Go, SQLite, dan Tailwind CSS.**

![Version](https://img.shields.io/badge/version-1.2-blue?style=flat-square)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat-square&logo=sqlite&logoColor=white)
![Tailwind](https://img.shields.io/badge/Tailwind_CSS-CDN-38BDF8?style=flat-square&logo=tailwindcss&logoColor=white)
![License](https://img.shields.io/badge/license-Private-red?style=flat-square)

*by **Baxium Coder***

</div>

---

## ✨ Fitur

| Modul | Keterangan |
|---|---|
| 👤 **Akun FB** | Kelola akun Facebook — email, password, cookie, 2FA, status |
| 📄 **Fanpage** | Manajemen fanpage beserta link langsung ke Facebook |
| 🏢 **Business Manager** | Data BM dan relasi ke akun FB |
| 📢 **Akun Iklan** | Ad account dengan limit harian/total, metode bayar |
| 🔵 **Pixel** | Tracking pixel dan relasi ke BM / akun FB |
| 📊 **Rekap / Summary** | Ringkasan keseluruhan data |

### 🔐 Keamanan
- Enkripsi AES-256-GCM untuk field sensitif (password, cookie, 2FA)
- Session-based authentication dengan durasi yang bisa dikonfigurasi
- CAPTCHA login (math / text)
- Proteksi CSRF

### 🎨 UI / UX
- Dark mode & Light mode
- Responsive — desktop dan mobile
- Live search / auto-filter real-time
- Modal detail info per record (read-only)
- Modal edit & riwayat per record
- Import bulk via TSV (copy-paste dari Google Sheets)
- Toast notification untuk aksi sukses/gagal

---

## 🚀 Instalasi & Menjalankan

### Prasyarat
- [Go 1.21+](https://golang.org/dl/)

### Clone & Jalankan

```bash
git clone https://github.com/username/dashboard-fb.git
cd dashboard-fb

# Install dependencies
go mod tidy

# Build binary
go build -o dashboard-fb.exe .

# Jalankan
./dashboard-fb.exe
```

Buka browser di `http://localhost:8080`

**Login default:**
```
Username : admin
Password : admin123
```

> ⚠️ Segera ganti password lewat menu **Pengaturan → Ganti Password** setelah login pertama.

---

## ⚙️ Konfigurasi

File `config.json` dibuat otomatis saat pertama kali dijalankan:

```json
{
  "username": "admin",
  "password_hash": "...",
  "encryption_key": "...",
  "port": "8080",
  "captcha_mode": "math",
  "session_hours": 1
}
```

| Key | Keterangan |
|---|---|
| `port` | Port server (default `8080`) |
| `encryption_key` | Kunci enkripsi AES-256 — **jangan dihapus/diganti** |
| `captcha_mode` | `math` atau `text` |
| `session_hours` | Durasi sesi login (jam) |

---

## 🗄️ Struktur Database

Database SQLite (`data.db`) dibuat otomatis. Tabel utama:

```
akun_fb       — Akun Facebook
fanpage       — Fanpage / halaman FB
bm            — Business Manager
akun_iklan    — Ad Account
pixel         — Facebook Pixel
riwayat       — Log riwayat per entitas
sessions      — Sesi login aktif
```

---

## 🚢 Deploy ke VPS

```bash
# 1. Push kode terbaru
git push origin main

# 2. Di VPS — pull & rebuild
git pull
go build -o dashboard-fb.exe .

# 3. Restart service
systemctl restart dashboard-fb   # atau pm2 restart / sesuai setup
```

> ⚠️ **Penting:** `config.json` dan `data.db` ada di `.gitignore` — **tidak ikut git**.
> Transfer kedua file ini manual ke VPS via SCP/SFTP dan **jangan sampai hilang**.
> Kehilangan `config.json` = data terenkripsi tidak bisa dibaca lagi.

```bash
# Transfer manual (hanya saat pertama deploy)
scp config.json user@vps:/path/to/dashboard-fb/
scp data.db     user@vps:/path/to/dashboard-fb/
```

---

## 📁 Struktur Proyek

```
dashboard-fb/
├── main.go              # Entry point
├── routes.go            # Routing
├── config.go            # Load / save config
├── crypto.go            # Enkripsi AES-256-GCM
├── funcmap.go           # Template helper functions
├── handlers/            # HTTP handler per modul
├── models/              # Query database per modul
├── db/
│   ├── db.go            # Inisialisasi SQLite
│   └── schema.sql       # Skema tabel
├── templates/
│   ├── layout.html      # Base layout + sidebar
│   ├── login.html       # Halaman login
│   ├── pages/           # Halaman per modul
│   └── partials/        # Fragment (riwayat, dll)
├── static/              # CSS, JS, assets
├── config.json          # ⚠️ Tidak di-git — buat otomatis
└── data.db              # ⚠️ Tidak di-git — database utama
```

---

## 📝 Changelog

### v1.2 — Mei 2026
- Tambah modal **Detail Info** di semua halaman (Akun FB, Fanpage, BM, Akun Iklan, Pixel)
- Info Akun FB: show/hide password, copy cookie & 2FA secret
- Fanpage: link langsung ke halaman Facebook
- **Live search** real-time dengan badge counter
- Dropdown BM di form Fanpage (opsional)
- Credits Baxium Coder di sidebar

### v1.1
- Redesign UI: dark/light mode, blue glassmorphism
- Riwayat per entitas
- Import bulk TSV

### v1.0
- Rilis awal — CRUD Akun FB, Fanpage, BM, Akun Iklan, Pixel

---

<div align="center">

Made with ❤️ by **Baxium Coder**

</div>
