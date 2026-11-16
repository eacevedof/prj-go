@echo off
rem install.bat -- lanzador de install.ps1 (mismo directorio) para instalar el siem-client en Windows con doble clic.
rem Pide permisos de administrador si no los tiene y pasa los argumentos tal cual:
rem   install.bat
rem   install.bat -Url https://siem.dominio.com -Token API_TOKEN
setlocal
set "SCRIPT=%~dp0install.ps1"
if not exist "%SCRIPT%" (
    echo ERROR: falta install.ps1 junto a este fichero. Descarga los dos desde el panel: Descargas.
    pause
    exit /b 1
)
net session >nul 2>&1
if errorlevel 1 (
    echo Pidiendo permisos de administrador...
    powershell -NoProfile -Command "Start-Process -Verb RunAs -FilePath powershell -ArgumentList '-NoProfile -ExecutionPolicy Bypass -NoExit -File \"%SCRIPT%\" %*'"
    exit /b 0
)
powershell -NoProfile -ExecutionPolicy Bypass -File "%SCRIPT%" %*
pause
