package ocr

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// ExtractText sends the image data to Gemini and returns the extracted text.
func ExtractText(imgData []byte) (string, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create genai client: %w", err)
	}
	defer client.Close()

	// Allow user to override model via environment variable
	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "gemini-3-flash-preview" // Default to newer model (2026)
	}

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0) // Deterministic output preferred for OCR

	prompt := []genai.Part{
		genai.ImageData("png", imgData),
		genai.Text("Extract all text from this image. Return ONLY the text, preserving layout if possible. Do not add markdown formatting or explanations."),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no text found in image")
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			result += string(txt)
		}
	}

	return result, nil
}
