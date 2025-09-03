@echo off
set app_name=task-mq
set target=a:\tmp

if not exist "%target%" mkdir "%target%"

go build -o "%target%\%app_name%.exe" .

if %ERRORLEVEL% neq 0 (
    echo compile fail!
    pause
    exit /b %ERRORLEVEL%
)

"%target%\%app_name%.exe"
