# Screenshot OCR Tool

A cross-platform (Windows currently optimized) tool that allows you to select a region of your screen, extracts the text using Google Gemini AI, and copies it to your clipboard.

## Features

- **Freeze-Screen Capture**: Captures your entire desktop (all monitors) and lets you drag to select.
- **AI-Powered OCR**: Uses Google's **Gemini 2.0 Flash** (configurable) for high-accuracy text extraction, including code blocks and complex layouts.
- **Clipboard Integration**: Automatically copies the result to your clipboard.
- **Stealth Mode**: Runs without a console window (on Windows) to avoid cluttering your screenshots.

## Prerequisites

1.  **Go**: Installed and in your PATH.
2.  **Gemini API Key**: Get one from [Google AI Studio](https://aistudio.google.com/).

## Setup

1.  Clone or download this repository.
2.  Set your API Key environment variable:
    ```powershell
    $env:GEMINI_API_KEY="your-api-key-here"
    ```
    *(Tip: Add this to your system environment variables to make it permanent.)*

3.  (Optional) Configure the Model:
    By default, it uses `gemini-2.0-flash`. You can override this:
    ```powershell
    $env:GEMINI_MODEL="gemini-1.5-pro"
    ```

## Building (Windows)

To build the application so it runs as a valid Windows GUI app (no command prompt window):

```powershell
.\build.bat
```

Or manually:
```powershell
go build -ldflags -H=windowsgui -o screenshot-ocr.exe .
```

## Usage

1.  Run `screenshot-ocr.exe`.
2.  Your screen will freeze with a dark overlay.
3.  **Click and Drag** to select the text you want to copy.
4.  Release the mouse button.
5.  Wait a moment for the AI to process.
6.  The text is now in your clipboard!

## Troubleshooting

-   **"GEMINI_API_KEY not set"**: Ensure you set the environment variable.
-   **Error Popup**: If the API fails (e.g., quota exceeded), a popup message will appear describing the error.
