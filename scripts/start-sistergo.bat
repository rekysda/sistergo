@echo off
REM Batch to build and run sistergo in new console
SET APPDIR=%~dp0\..\
PUSHD %APPDIR%
IF NOT EXIST "%APPDIR%cmd\server\main.go" (
    echo Cannot find main.go under %APPDIR%cmd\server
    popd
    exit /b 1
)

REM Build
go build -o sistergo.exe ./cmd/server
IF NOT EXIST sistergo.exe (
    echo Build failed
    popd
    exit /b 1
)

start "SisterGo" cmd /k "%cd%\sistergo.exe"
popd
