@echo off
title FB Manager - Build & Run
echo.
echo  ================================
echo   FB Manager - Build and Run
echo  ================================
echo.

echo  [1/3] Stop server lama...
taskkill /F /IM dashboard-fb.exe /T >nul 2>&1
timeout /t 1 /nobreak >nul

echo  [2/3] Build binary...
go build -v -o dashboard-fb.exe .
if %errorlevel% neq 0 (
    echo  GAGAL: Build error!
    pause
    exit /b 1
)

echo  [3/3] Build selesai! Menjalankan server...
echo.
echo  ================================
echo   Buka browser: http://localhost:8080
echo   Tekan Ctrl+C untuk stop.
echo  ================================
echo.
dashboard-fb.exe
pause
