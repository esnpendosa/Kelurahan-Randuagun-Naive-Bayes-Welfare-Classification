@echo off
set APP_NAME=Klasifikasi-Warga-Randuagung
echo ========================================================
echo   Membangun %APP_NAME% (.exe)
echo ========================================================
echo.
echo Sedang mengompilasi... Mohon tunggu...
echo.

:: Mengompilasi dengan flag -H=windowsgui untuk menyembunyikan jendela CMD
:: Flag -s -w digunakan untuk memperkecil ukuran file binary
go build -ldflags "-H=windowsgui -s -w" -o "%APP_NAME%.exe" main.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo [SUKSES] Aplikasi berhasil dibangun: %APP_NAME%.exe
    echo.
    echo Anda sekarang dapat menjalankan %APP_NAME%.exe.
    echo Browser akan otomatis terbuka saat aplikasi dijalankan.
    echo.
) else (
    echo.
    echo [GAGAL] Terjadi kesalahan saat membangun aplikasi.
    echo Pastikan Go sudah terinstal dan semua dependensi sudah ada.
    echo.
)

pause
