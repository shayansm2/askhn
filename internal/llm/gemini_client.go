package llm

import (
	"context"
	"log"

	"google.golang.org/genai"
)

type GeminiLLM struct {
	client *genai.Client
	model  string
	ctx    context.Context
}

func NewGeminiLLM(apiKey, model string) GeminiLLM {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("Error creating Gemini client: %v", err)
	}
	return GeminiLLM{
		client: client,
		model:  model,
		ctx:    ctx,
	}
}

func (llm *GeminiLLM) Chat(message Message) (string, error) {
	var config *genai.GenerateContentConfig
	if message.System != "" {
		systemInstruction := genai.NewContentFromText(message.System, genai.RoleUser)
		config = &genai.GenerateContentConfig{
			SystemInstruction: systemInstruction,
		}
	}

	result, err := llm.client.Models.GenerateContent(
		llm.ctx,
		llm.model,
		genai.Text(message.User),
		config,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}
