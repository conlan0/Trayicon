@echo off
set "targetDir=C:\Program Files\TacticalAgent"
set "exePath=%targetDir%\tray.exe"

:: Check if the tray icon executable exists and delete it if it does
if exist "%exePath%" (
    del "%exePath%"
)

cd /d ..\..\..
mkdir "install"
cd "install"

curl -o "tray.exe" https://example.com/downloads/tray.exe --ssl-no-revoke
copy "tray.exe" "%targetDir%"

reg add "HKLM\Software\Microsoft\Windows\CurrentVersion\Run" /v "TacticalRMMTray" /t REG_SZ /d "%exePath%" /f
