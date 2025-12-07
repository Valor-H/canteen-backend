@echo off
echo Starting Canteen Backend System...

echo Starting Redis Server...
start "Redis Server" /min D:\Tools\Redis\redis-server.exe

echo Starting MySQL Server...
start "MySQL Server" /min D:\Tools\MySQL\bin\mysqld --defaults-file="D:\Tools\MySQL\my.ini"

echo Waiting for services to start...
timeout /t 3 /nobreak >nul

echo Starting Go Application...
cd /d d:\Codes\canteen-backend
go run ./cmd/server/main.go

pause