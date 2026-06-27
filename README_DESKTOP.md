# Aplikasi Desktop - Troubleshooting Guide

## ✅ Status Build
- Chrome/Edge: **TERDETEKSI** ✓
- File .exe: **SUDAH DIBUAT** ✓
- Build Time: 27/6/2026 11:18 AM

## 🔧 Dua Versi File

### 1. `Klasifikasi-Warga-Randuagung.exe` (Versi Production)
- Flag: `-H=windowsgui` (tidak show console)
- Ukuran: ~18 MB
- Untuk distribusi ke user

### 2. `Klasifikasi-Warga-Randuagung-debug.exe` (Versi Debug)
- Flag: Tanpa `-H=windowsgui` (show console)
- Ukuran: ~18 MB
- **GUNAKAN INI UNTUK TEST** - bisa lihat error

## 🚀 Cara Test (PENTING!)

### Gunakan Versi Debug untuk Test:

**Double-click:** `run-debug.bat`

Ini akan:
1. Jalankan versi debug
2. Show console output
3. Tampilkan error jika ada

### Apa yang Seharusnya Terjadi:

#### Skenario A: Desktop Window Berhasil (TUJUAN KITA)
```
Server berjalan di http://127.0.0.1:8080
✅ Desktop window berhasil dibuat!
⏹️  Tutup window untuk menghentikan aplikasi.
```
→ Desktop window muncul tanpa address bar

#### Skenario B: Desktop Window Gagal (Fallback ke Browser)
```
Server berjalan di http://127.0.0.1:8080
⚠️  Gagal membuat desktop window: [error detail]
📱 Membuka di browser sebagai fallback...
🌐 URL: http://127.0.0.1:8080
✅ Server tetap berjalan. Buka browser dan akses URL di atas.
```
→ Server tetap jalan, buka browser manual ke http://localhost:8080

## ❓ Troubleshooting

### Masalah 1: Aplikasi Langsung Berhenti
**Penyebab:** Lorca tidak bisa start Chrome

**Solusi:**
1. Gunakan `run-debug.bat` untuk lihat error detail
2. Jika fallback ke browser, buka http://localhost:8080 manual
3. Screenshot console output dan kasih tau saya

### Masalah 2: Error "Chrome not found"
**Penyebab:** Lorca tidak menemukan Chrome executable

**Solusi A - Set LORCA_CHROME_PATH:**
```batch
set LORCA_CHROME_PATH=C:\Program Files\Google\Chrome\Application\chrome.exe
run-debug.bat
```

**Solusi B - Reinstall Chrome:**
1. Uninstall Chrome
2. Download fresh: https://www.google.com/chrome/
3. Install ke default location
4. Test ulang

### Masalah 3: Port 8080 Sudah Dipakai
**Penyebab:** Aplikasi lain pakai port 8080

**Cek port:**
```cmd
netstat -ano | findstr :8080
```

**Solusi:**
1. Tutup aplikasi yang pakai port 8080 (mungkin Laragon/Apache)
2. Atau edit main.go ubah port 8080 ke 8081

### Masalah 4: Database Error
**Penyebab:** data_skripsi.db corrupt atau permission issue

**Solusi:**
1. Backup data_skripsi.db
2. Rename atau delete data_skripsi.db
3. Aplikasi akan auto-create database baru
4. Test ulang

## 📊 Log Output yang Bagus

Jika ada masalah, saya butuh lihat:

1. **Console output lengkap** dari run-debug.bat
2. **Error message** jika ada
3. **Task Manager** - cek apakah process `chrome.exe` atau `msedge.exe` muncul

## 🎯 Next Steps

### Jika Desktop Window Berhasil:
✅ Build ulang production version:
```batch
rebuild.bat
```
✅ Distribusikan `Klasifikasi-Warga-Randuagung.exe`

### Jika Fallback ke Browser:
Ada 2 opsi:

**Opsi 1: Fix Lorca Issue**
- Debug kenapa Lorca gagal
- Fix dan build ulang

**Opsi 2: Pakai Browser (Solusi Sementara)**
- Server tetap jalan
- User buka http://localhost:8080 manual
- Masih lebih baik dari sebelumnya (tidak auto-open browser)

## 🔍 Technical Details

### Lorca Requirements:
- Chrome atau Edge installed
- Chrome executable di lokasi default
- Windows 10/11
- No sandboxing/restricted environment

### Fallback Behavior:
- Jika Lorca gagal, server tetap jalan
- User bisa akses via browser manual
- Tidak auto-open browser (lebih baik dari versi lama)

### Build Flags:
```bash
# Production (no console)
go build -ldflags "-H=windowsgui -s -w" -o Klasifikasi-Warga-Randuagung.exe main.go

# Debug (with console)
go build -ldflags "-s -w" -o Klasifikasi-Warga-Randuagung-debug.exe main.go
```

## 📞 Report Bug

Jika masih gagal, screenshot:
1. Output dari `run-debug.bat`
2. Task Manager (cari chrome.exe)
3. Hasil dari `check-chrome.bat`

Kirim ke saya untuk analisa lebih lanjut.

---

**Last Updated:** 27 Juni 2026, 11:18 AM
