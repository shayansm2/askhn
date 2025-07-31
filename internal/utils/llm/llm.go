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

type Message struct {
	User   string
	System string
	Err    error
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
