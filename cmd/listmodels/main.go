package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	iter := client.ListModels(ctx)
	fmt.Println("Available Models:")
	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error listing models: %v", err)
		}
		fmt.Printf("- %s (%s)\n", m.Name, m.DisplayName)
	}
}
