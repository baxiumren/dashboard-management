@echo off
title FB Manager
echo.
echo  ================================
echo   FB Manager - Ads Dashboard
echo  ================================
echo.

if not exist dashboard-fb.exe (
    echo  File dashboard-fb.exe belum ada.
    echo  Jalankan build.bat dulu!
    echo.
    pause
    exit /b 1
)

echo  Menjalankan server dari binary...
echo  Buka browser: http://localhost:8080
echo.
echo  Tekan Ctrl+C untuk stop.
echo.
dashboard-fb.exe
pause
