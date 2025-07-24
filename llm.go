package temporallm

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type LargeLanguageModel struct {
	client openai.Client
	model  string
}

func NewLLM(baseURL, apiKey, model string) LargeLanguageModel {
	clientOptions := []option.RequestOption{}
	if baseURL != "" {
		clientOptions = append(clientOptions, option.WithBaseURL(baseURL))
	}
	if apiKey != "" {
		clientOptions = append(clientOptions, option.WithAPIKey(apiKey))
	}
	return LargeLanguageModel{
		client: openai.NewClient(clientOptions...),
		model:  model,
	}
}

func (llm *LargeLanguageModel) Chat(message string) (string, error) {
	chatCompletion, err := llm.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
		Model: llm.model,
	})
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
