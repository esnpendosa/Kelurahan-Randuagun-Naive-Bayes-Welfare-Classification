@echo off
echo ========================================
echo   CEK CHROME / EDGE
echo ========================================
echo.

echo Mencari Google Chrome...
if exist "C:\Program Files\Google\Chrome\Application\chrome.exe" (
    echo [OK] Google Chrome ditemukan!
    echo Lokasi: C:\Program Files\Google\Chrome\Application\chrome.exe
) else if exist "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" (
    echo [OK] Google Chrome ditemukan!
    echo Lokasi: C:\Program Files (x86)\Google\Chrome\Application\chrome.exe
) else (
    echo [X] Google Chrome TIDAK ditemukan!
)
echo.

echo Mencari Microsoft Edge...
if exist "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" (
    echo [OK] Microsoft Edge ditemukan!
    echo Lokasi: C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe
) else if exist "C:\Program Files\Microsoft\Edge\Application\msedge.exe" (
    echo [OK] Microsoft Edge ditemukan!
    echo Lokasi: C:\Program Files\Microsoft\Edge\Application\msedge.exe
) else (
    echo [X] Microsoft Edge TIDAK ditemukan!
)
echo.

echo ========================================
echo   DIAGNOSIS
echo ========================================
echo.
echo Lorca (library desktop app) membutuhkan:
echo - Google Chrome (recommended) ATAU
echo - Microsoft Edge
echo.
echo Jika keduanya tidak ditemukan, aplikasi tidak akan bisa jalan.
echo.
echo Solusi:
echo 1. Install Google Chrome dari: https://www.google.com/chrome/
echo 2. Atau update Microsoft Edge
echo.

pause
