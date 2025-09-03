@echo off
set target=a:\tmp

if not exist "%target%" mkdir "%target%"

go build -o "%target%\user-api.exe" .

if %ERRORLEVEL% neq 0 (
    echo compile fail!
    pause
    exit /b %ERRORLEVEL%
)

"%target%\user-api.exe"
