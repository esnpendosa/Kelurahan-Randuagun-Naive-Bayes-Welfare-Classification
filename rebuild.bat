@echo off
echo ========================================================
echo   REBUILD Aplikasi Desktop
echo ========================================================
echo.

echo [1/3] Menutup aplikasi yang running...
taskkill /F /IM Klasifikasi-Warga-Randuagung.exe 2>nul
if %ERRORLEVEL% EQU 0 (
    echo     Aplikasi berhasil ditutup.
) else (
    echo     Tidak ada aplikasi yang running.
)
echo.

echo [2/3] Menghapus file exe lama...
if exist Klasifikasi-Warga-Randuagung.exe (
    del Klasifikasi-Warga-Randuagung.exe
    echo     File lama berhasil dihapus.
) else (
    echo     File exe tidak ditemukan.
)
echo.

echo [3/3] Building aplikasi baru...
go build -ldflags "-H=windowsgui -s -w" -o Klasifikasi-Warga-Randuagung.exe main.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================================
    echo   ✅ REBUILD BERHASIL!
    echo ========================================================
    echo.
    echo File: Klasifikasi-Warga-Randuagung.exe
    echo.
    echo Double-click file .exe untuk menjalankan aplikasi.
    echo.
) else (
    echo.
    echo ========================================================
    echo   ❌ REBUILD GAGAL!
    echo ========================================================
    echo.
    echo Pastikan:
    echo - Anda berada di folder project yang benar
    echo - main.go ada di folder ini
    echo - go.mod sudah di-update
    echo.
)

pause
