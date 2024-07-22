package gemini

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

const (
	location  = "northamerica-northeast1"
	modelName = "gemini-1.5-pro"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(ctx context.Context) (*GeminiClient, error) {
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT environment variable not set")
	}

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &GeminiClient{client: client}, nil
}
