package llm

import (
	"context"
	"encoding/json"

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

func NewOpenAiLLM(baseURL, apiKey, model string) LargeLanguageModel {
	return LargeLanguageModel{
		client: openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(apiKey),
		),
		model: model,
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

// todo
func (llm *LargeLanguageModel) ChatWithJSONFormat(message Message, result interface{}) error {
	response, err := llm.Chat(message)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return err
	}
	return nil
}

// todo
func (llm *LargeLanguageModel) chatWithCriteria(message Message, check func(string) error) (string, error) {
	response, err := llm.Chat(message)
	if err != nil {
		return "", err
	}
	err = check(response)
	if err != nil {
		return "", err
	}
	return response, nil
}
