@echo off

set "DEST=%~dp0"
set "SRC=C:\Users\byte\repos\6-public\abfallkalender-api"

echo Compiling CLI for multiple platforms...
echo Source: %SRC%
echo Output: %DEST%bin

if not exist "%DEST%bin" mkdir "%DEST%bin"

pushd "%SRC%"
if errorlevel 1 (
    echo ERROR: source dir not found: %SRC%
    exit /b 1
)

echo.
echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -o "%DEST%bin\abfallkalender.exe" ./cmd
if errorlevel 1 (
    popd
    echo ERROR: Windows build failed
    exit /b 1
)

echo.
echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -o "%DEST%bin\abfallkalender" ./cmd
if errorlevel 1 (
    popd
    echo ERROR: Linux build failed
    exit /b 1
)

echo.
echo Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -o "%DEST%bin\abfallkalender-darwin" ./cmd
if errorlevel 1 (
    popd
    echo ERROR: macOS build failed
    exit /b 1
)

popd

echo.
echo Compilation complete.
echo.
echo SHA256 Hashes:
echo.
certutil -hashfile "%DEST%bin\abfallkalender.exe" SHA256 | findstr /v "hash"
certutil -hashfile "%DEST%bin\abfallkalender" SHA256 | findstr /v "hash"
certutil -hashfile "%DEST%bin\abfallkalender-darwin" SHA256 | findstr /v "hash"

pause
