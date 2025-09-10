@echo off
set target=a:\tmp

if not exist "%target%" mkdir "%target%"

go build -o "%target%\im-rpc.exe" .

if %ERRORLEVEL% neq 0 (
    echo compile fail!
    pause
    exit /b %ERRORLEVEL%
)

"%target%\im-ws.exe"
