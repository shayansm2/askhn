package chatbot

import (
	"context"

	"github.com/shayansm2/temporallm/internal/utils/llm"
	"github.com/shayansm2/temporallm/internal/utils/textdb"
)

type LLMActivities struct {
	LLM *llm.LargeLanguageModel
}

func (i *LLMActivities) SimpleChat(ctx context.Context, message string) (string, error) {
	return i.LLM.Chat(llm.Message{User: message})
}

type HackerNewsApiActivities struct{}

func (i *HackerNewsApiActivities) RetrieveHackerNewsItem(ctx context.Context, id int) (*HackerNewsResponse, error) {
	return HackerNewsItem(id)
}

type IndexRequest struct {
	IndexName string
	Document  *HackerNewsResponse
}

func (i *HackerNewsApiActivities) IndexInElasticsearch(ctx context.Context, req IndexRequest) error {
	return textdb.Index(req.IndexName, req.Document)
}
