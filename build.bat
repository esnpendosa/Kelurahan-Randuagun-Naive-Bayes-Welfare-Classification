@echo off
set APP_NAME=Klasifikasi-Warga-Randuagung
echo ========================================================
echo   Membangun %APP_NAME% (.exe)
echo ========================================================
echo.
echo Sedang mengompilasi... Mohon tunggu...
echo.

:: Download dependencies terlebih dahulu
echo Mengunduh dependencies...
go mod download
go mod tidy
echo.

:: Mengompilasi dengan flag -H=windowsgui untuk menyembunyikan jendela CMD
:: Flag -s -w digunakan untuk memperkecil ukuran file binary
echo Mengompilasi aplikasi desktop...
go build -ldflags "-H=windowsgui -s -w" -o "%APP_NAME%.exe" main.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo [SUKSES] Aplikasi berhasil dibangun: %APP_NAME%.exe
    echo.
    echo Aplikasi Desktop Windows Native siap digunakan!
    echo Aplikasi akan berjalan dalam window sendiri (TIDAK di browser).
    echo.
    echo Menggunakan Lorca (Chrome embedded window)
    echo.
    echo Double-click %APP_NAME%.exe untuk menjalankan.
    echo.
) else (
    echo.
    echo [GAGAL] Terjadi kesalahan saat membangun aplikasi.
    echo.
    echo Pastikan:
    echo 1. Go sudah terinstal (go version ^>= 1.18)
    echo 2. Jalankan: go mod tidy
    echo.
    echo Jika error "Chrome not found" saat run:
    echo - Install Google Chrome dari: https://www.google.com/chrome/
    echo.
)

pause
