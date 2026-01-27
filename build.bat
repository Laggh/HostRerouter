@echo off
go build -ldflags "-H=windowsgui" -o HostRerouter.exe .
echo Build complete for HostRerouter.exe
pause