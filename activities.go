package temporallm

import (
	"context"
)

type LLMActivities struct {
	LLM *LargeLanguageModel
}

func (i *LLMActivities) Chat(ctx context.Context, message string) (string, error) {
	return i.LLM.Chat(message)
}
