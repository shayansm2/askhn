package chatbot

import (
	"context"

	"github.com/shayansm2/temporallm/internal/utils/llm"
)

type LLMActivities struct {
	LLM *llm.LargeLanguageModel
}

func (i *LLMActivities) SimpleChat(ctx context.Context, message string) (string, error) {
	return i.LLM.Chat(llm.Message{User: message})
}
