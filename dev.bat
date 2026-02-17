@echo off
echo Starting Social Media Lead Automation (Local Dev)...

:: Check if Redis is running (optional but recommended)
docker ps | findstr "redis" >nul
if %errorlevel% neq 0 (
    echo [WARNING] Redis is not running in Docker. Some features may be disabled.
    echo To start Redis: docker run -d -p 6379:6379 redis
)

:: Start Backend
start "LeadPilot Backend" cmd /k "cd backend && go run ./cmd/api"

:: Start Frontend
timeout /t 10 /nobreak >nul
start "LeadPilot Frontend" cmd /k "cd frontend && npm run dev"

echo.
echo Services starting...
echo Backend: http://localhost:8080
echo Frontend: http://localhost:3000
echo.
echo Press any key to close this launcher (services will keep running)...
pause >nul
