@echo off
setlocal enabledelayedexpansion

:: Default port is 8080, but can be overridden with PORT environment variable
if "%PORT%"=="" set PORT=8080

echo Checking if port %PORT% is in use...

:: Find process using the port
for /f "tokens=5" %%a in ('netstat -ano ^| findstr /r ":%PORT% "') do (
    set PID=%%a
    goto :found
)

:notfound
echo Port %PORT% is not in use.
goto :startserver

:found
echo Port %PORT% is in use by process with PID %PID%. Killing process...
taskkill /F /PID %PID%
if %ERRORLEVEL% EQU 0 (
    echo Process killed.
    :: Wait a moment to ensure the port is released
    timeout /t 1 /nobreak >nul
) else (
    echo Failed to kill process. You may need to run this script as administrator.
)

:startserver
echo Starting server on port %PORT%...
set PORT=%PORT%
server.exe

endlocal
