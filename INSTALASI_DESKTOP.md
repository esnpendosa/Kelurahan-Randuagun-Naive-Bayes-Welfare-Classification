# Instalasi Aplikasi Desktop Windows

## Persyaratan Sistem

1. **Windows 10 atau Windows 11**
2. **Go 1.18 atau lebih baru** (sudah terinstall)
3. **GCC/MinGW** (diperlukan untuk kompilasi webview)

## Cara Install GCC/MinGW di Windows

Webview memerlukan GCC untuk kompilasi. Berikut cara instalasi:

### Metode 1: TDM-GCC (Recommended - Paling Mudah)

1. Download TDM-GCC dari: https://github.com/jmeubank/tdm-gcc/releases
2. Pilih file installer terbaru (contoh: `tdm-gcc-10.3.0.exe`)
3. Install dengan pengaturan default
4. Restart terminal/command prompt
5. Test dengan: `gcc --version`

### Metode 2: MSYS2

1. Download MSYS2 dari: https://www.msys2.org/
2. Install MSYS2
3. Buka MSYS2 terminal dan jalankan:
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```
4. Tambahkan ke PATH: `C:\msys64\mingw64\bin`
5. Test dengan: `gcc --version`

### Metode 3: WinLibs

1. Download dari: https://winlibs.com/
2. Extract file ZIP ke folder (contoh: `C:\mingw64`)
3. Tambahkan ke PATH: `C:\mingw64\bin`
4. Restart terminal
5. Test dengan: `gcc --version`

## Build Aplikasi Desktop

Setelah GCC terinstall:

```bash
# 1. Download dependencies
go mod download
go mod tidy

# 2. Build aplikasi
build.bat

# Atau manual:
go build -ldflags "-H=windowsgui -s -w" -o Klasifikasi-Warga-Randuagung.exe main.go
```

## Menjalankan Aplikasi

1. Double-click file `Klasifikasi-Warga-Randuagung.exe`
2. Aplikasi akan membuka window sendiri (BUKAN browser)
3. Window akan menampilkan antarmuka aplikasi

## Troubleshooting

### Error: "gcc: command not found"

**Solusi:** GCC belum terinstall atau belum ada di PATH. Install GCC menggunakan salah satu metode di atas.

### Error: "exec: gcc: executable file not found in %PATH%"

**Solusi:** 
1. Pastikan GCC sudah terinstall: `gcc --version`
2. Jika sudah install, restart terminal
3. Jika masih error, tambahkan path GCC ke System Environment Variables

### Aplikasi masih membuka browser

**Solusi:** Pastikan Anda build ulang setelah perubahan kode:
```bash
# Hapus file exe lama
del Klasifikasi-Warga-Randuagung.exe

# Build ulang
build.bat
```

### Window terlalu kecil/besar

Edit `main.go` baris:
```go
w.SetSize(1280, 800, webview.HintNone)
```

Ubah 1280x800 sesuai kebutuhan.

## Distribusi Aplikasi

Setelah build berhasil, Anda bisa distribusikan file:

```
Klasifikasi-Warga-Randuagung.exe  (aplikasi utama)
data_skripsi.db                    (database, akan auto-create jika tidak ada)
```

**Catatan:** User tidak perlu install Go atau GCC untuk menjalankan file .exe yang sudah di-build.

## Fitur Aplikasi Desktop

✅ Window aplikasi native Windows
✅ Tidak membuka browser eksternal
✅ Bisa minimize, maximize, close seperti aplikasi normal
✅ Ukuran window: 1280x800 (bisa disesuaikan)
✅ Title bar: "Sistem Klasifikasi Kesejahteraan Randuagung"
✅ Single executable file
✅ Database SQLite terintegrasi

## Update dari Versi Browser ke Desktop

Perubahan yang dilakukan:
1. Mengganti `os/exec` (buka browser) dengan `webview` (window native)
2. Server tetap berjalan di background (localhost)
3. UI ditampilkan di webview window, bukan browser
4. Semua fitur existing tetap berfungsi normal

## Kontak

Jika ada masalah, hubungi developer atau buat issue di repository.
