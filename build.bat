@echo off
echo Building Screenshot OCR Tool...
"C:\Program Files\Go\bin\go.exe" build -ldflags -H=windowsgui -o screenshot-ocr.exe .
if %errorlevel% neq 0 (
    echo Build failed!
    pause
    exit /b %errorlevel%
)
echo Build successful: screenshot-ocr.exe
