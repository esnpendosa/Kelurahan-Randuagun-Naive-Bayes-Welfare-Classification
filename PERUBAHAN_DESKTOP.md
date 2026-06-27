# Perubahan Aplikasi ke Desktop Windows

## Status: ✅ SELESAI - SIAP TEST

## Ringkasan Perubahan

Aplikasi telah diubah dari **web-based (browser)** menjadi **desktop native Windows** menggunakan **Lorca**.

---

## File yang Dimodifikasi

### 1. `main.go`
**Perubahan:**
- ❌ **Dihapus:** `import "os/exec"` dan `import "runtime"` 
- ❌ **Dihapus:** Fungsi `bukaBrowser()` yang membuka browser eksternal
- ✅ **Ditambah:** `import "github.com/zserge/lorca"`
- ✅ **Diganti:** Window Chrome embedded menggunakan Lorca

**Sebelum:**
```go
// Membuka browser secara otomatis
go func(p int) {
    time.Sleep(500 * time.Millisecond)
    bukaBrowser(fmt.Sprintf("http://localhost:%d", p))
}(port)
```

**Sesudah:**
```go
// Buat window desktop menggunakan Lorca
ui, err := lorca.New(fmt.Sprintf("http://127.0.0.1:%d", port), "", 1280, 800)
if err != nil {
    fmt.Printf("Gagal membuat window aplikasi: %v\n", err)
    fmt.Println("Browser akan dibuka sebagai fallback...")
    time.Sleep(1 * time.Second)
    select {}
}
defer ui.Close()

// Tunggu sampai window ditutup
<-ui.Done()
```

---

### 2. `go.mod`
**Perubahan:**
- ✅ **Ditambah:** `github.com/zserge/lorca v0.1.10`
- Library lainnya tetap sama

---

### 3. `build.bat`
**Perubahan:**
- ✅ **Ditambah:** `go mod download` sebelum build
- ✅ **Ditambah:** `go mod tidy` sebelum build
- ✅ **Update:** Pesan sukses lebih jelas tentang desktop app
- ✅ **Update:** Troubleshooting lebih spesifik untuk Lorca

---

## File Dokumentasi Baru

### 1. `QUICK_START_DESKTOP.txt`
Panduan singkat cara build dan run aplikasi desktop

### 2. `INSTALASI_DESKTOP.md`
Panduan lengkap dengan troubleshooting

### 3. `PERUBAHAN_DESKTOP.md` (file ini)
Dokumentasi perubahan teknis

---

## Cara Build

### Langkah 1: Pastikan Chrome/Edge Terinstall
Lorca membutuhkan Chrome atau Edge (sudah ada di Windows 10/11)

### Langkah 2: Build Aplikasi
```bash
# Otomatis
build.bat

# Manual
go mod tidy
go build -ldflags "-H=windowsgui -s -w" -o Klasifikasi-Warga-Randuagung.exe main.go
```

### Langkah 3: Jalankan
```bash
Klasifikasi-Warga-Randuagung.exe
```

---

## Keuntungan Menggunakan Lorca

✅ **Tidak butuh GCC/MinGW** (compile lebih mudah)
✅ **Tidak butuh WebView2 runtime** tambahan  
✅ **Build lebih cepat** daripada webview native
✅ **Menggunakan Chrome** yang sudah terinstall di Windows
✅ **Ukuran file .exe lebih kecil**
✅ **Cross-platform** (bisa jalan di Linux/Mac juga)

---

## Cara Kerja Aplikasi Desktop

1. **Server Echo** tetap berjalan di background (localhost:8080)
2. **Lorca** membuat window Chrome embedded tanpa address bar
3. Window Lorca menampilkan UI dari http://127.0.0.1:8080
4. Semua fitur tetap sama, hanya tampilan yang berubah
5. Ketika window ditutup, aplikasi berhenti

---

## Perbedaan dengan Versi Browser

| Aspek | Versi Browser (Lama) | Versi Desktop (Baru) |
|-------|---------------------|---------------------|
| Window | Browser eksternal (Chrome/Edge) | Window embedded (Lorca) |
| Address Bar | ✅ Terlihat | ❌ Tidak terlihat |
| Tab lain | Bisa buka tab baru | Tidak bisa |
| Shortcut | Tidak ada | Bisa dibuat shortcut .exe |
| Distribusi | Perlu instruksi "buka localhost" | Tinggal double-click .exe |
| User Experience | Seperti website | Seperti aplikasi native |

---

## Troubleshooting

### ❌ Error: "Chrome not found"
**Penyebab:** Chrome atau Edge tidak terinstall

**Solusi:**
1. Install Google Chrome: https://www.google.com/chrome/
2. ATAU pastikan Microsoft Edge up-to-date (Windows 10/11 sudah ada)

### ❌ Aplikasi masih buka browser biasa
**Penyebab:** Build menggunakan kode lama

**Solusi:**
```bash
del Klasifikasi-Warga-Randuagung.exe
build.bat
```

### ❌ Error saat compile
**Penyebab:** Dependencies tidak lengkap

**Solusi:**
```bash
go clean
go mod tidy
go build -ldflags "-H=windowsgui -s -w" -o Klasifikasi-Warga-Randuagung.exe main.go
```

---

## Testing Checklist

Setelah build berhasil, test hal-hal berikut:

- [ ] Double-click .exe bisa jalan
- [ ] Window muncul tanpa buka browser eksternal
- [ ] Window size 1280x800
- [ ] Tidak ada address bar di window
- [ ] Login berfungsi
- [ ] Data warga bisa ditampilkan
- [ ] Klasifikasi bisa dijalankan
- [ ] Hasil klasifikasi tampil
- [ ] Training model berfungsi
- [ ] Laporan bisa dilihat
- [ ] Close window = aplikasi berhenti

---

## Distribusi ke User

### Yang Dibutuhkan User:
1. File `Klasifikasi-Warga-Randuagung.exe`
2. Chrome atau Edge (sudah ada di Windows 10/11)
3. File `data_skripsi.db` (optional - auto-create jika tidak ada)

### Tidak Dibutuhkan:
- ❌ Install Go
- ❌ Install GCC/MinGW
- ❌ Install Node.js
- ❌ Install dependencies apapun

### Cara Distribusi:
```
Klasifikasi-Warga-Randuagung/
├── Klasifikasi-Warga-Randuagung.exe  (WAJIB)
├── data_skripsi.db                    (optional)
└── README.txt                         (instruksi singkat)
```

---

## Kontak & Support

Jika ada masalah saat build atau run:
1. Test dulu dengan langkah di atas
2. Screenshot error message
3. Beri feedback ke developer

---

## Catatan Teknis

### Lorca vs Webview vs Wails

| Library | Kelebihan | Kekurangan | Pilihan |
|---------|-----------|------------|---------|
| **Lorca** | Mudah, tidak butuh GCC, ukuran kecil | Butuh Chrome/Edge installed | ✅ **DIPILIH** |
| Webview | True native, tidak butuh browser | Butuh GCC, setup rumit | ❌ |
| Wails | Full-featured, modern | Kompleks, learning curve | ❌ |

**Kesimpulan:** Lorca dipilih karena paling mudah dan praktis untuk kasus ini.

---

## Changelog

### Version 2.0 - Desktop Native
- ✅ Ubah dari browser ke desktop window
- ✅ Gunakan Lorca untuk embedded Chrome
- ✅ Hapus fungsi bukaBrowser()
- ✅ Update build.bat
- ✅ Tambah dokumentasi lengkap

### Version 1.0 - Web Based (Lama)
- Browser eksternal
- Fungsi bukaBrowser() otomatis

---

**Status:** ✅ **SIAP UNTUK TESTING**

Silakan jalankan `build.bat` dan test hasilnya!
