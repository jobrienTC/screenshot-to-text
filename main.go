package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"

	"screenshot-to-text/internal/capture"
	"screenshot-to-text/internal/ocr"
	"screenshot-to-text/internal/ui"

	"github.com/atotto/clipboard"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/sys/windows"
)

func main() {
	// Setup panic/error handling to show a message box if we crash,
	// because in GUI mode (no console) the user won't see log.Fatal output.
	defer func() {
		if r := recover(); r != nil {
			showError(fmt.Sprintf("Panic: %v", r))
		}
	}()

	if os.Getenv("GEMINI_API_KEY") == "" {
		showError("GEMINI_API_KEY environment variable is not set")
		return
	}

	// 1. Capture Screen
	fullScreenshot, err := capture.GetScreenShot()
	if err != nil {
		showError(fmt.Sprintf("Failed to capture screen: %v", err))
		return
	}

	// 2. Initialize UI with screenshot
	// Ebiten settings for fullscreen overlay
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true) // Always on top
	// Try to maximize/fullscreen to cover everything
	ebiten.SetFullscreen(true)

	selectionApp := ui.NewSelectionApp(fullScreenshot)

	if err := ebiten.RunGame(selectionApp); err != nil {
		if err == ebiten.Termination {
			// Normal exit
		} else {
			showError(fmt.Sprintf("UI Error: %v", err))
			return
		}
	}

	// 3. Get Selection
	selectedImg := selectionApp.Result()
	if selectedImg == nil {
		// No selection made or cancelled.
		return
	}

	// 4. Convert to PNG bytes for API
	var buf bytes.Buffer
	if err := png.Encode(&buf, selectedImg); err != nil {
		showError(fmt.Sprintf("Failed to encode image: %v", err))
		return
	}

	// 5. Run OCR
	text, err := ocr.ExtractText(buf.Bytes())
	if err != nil {
		showError(fmt.Sprintf("OCR Failed: %v", err))
		return
	}

	if text == "" {
		showError("No text detected in the selection.")
		return
	}

	// 6. Copy to Clipboard
	if err := clipboard.WriteAll(text); err != nil {
		showError(fmt.Sprintf("Failed to copy to clipboard: %v", err))
		return
	}
}

// showError displays a Windows MessageBox.
// This is critical for the "no console" build.
func showError(msg string) {
	// Fallback to console log just in case
	log.Println(msg)

	caption := windows.StringToUTF16Ptr("Screenshot OCR Error")
	text := windows.StringToUTF16Ptr(msg)

	// MessageBoxW(hwnd, text, caption, type)
	// Type 0x10 = MB_ICONERROR
	windows.MessageBox(0, text, caption, 0x10)
}
