package llm

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type LargeLanguageModel struct {
	client openai.Client
	model  string
}

func NewOllamaLLM(baseURL, model string) LargeLanguageModel {
	return LargeLanguageModel{
		client: openai.NewClient(option.WithBaseURL(baseURL)),
		model:  model,
	}
}

func NewOpenAiLLM(apiKey, model string) LargeLanguageModel {
	return LargeLanguageModel{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
		model:  model,
	}
}

func (llm *LargeLanguageModel) Chat(message Message) (string, error) {
	chatCompletion, err := llm.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message.User),
			openai.SystemMessage(message.System),
		},
		Model: llm.model,
	})
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
