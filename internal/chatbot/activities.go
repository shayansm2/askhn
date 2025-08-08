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

func (i *LLMActivities) Chat(ctx context.Context, userMsg, systemMsg string) (string, error) {
	return i.LLM.Chat(llm.Message{User: userMsg, System: systemMsg})
}

type HackerNewsApiActivities struct{}

func (i *HackerNewsApiActivities) RetrieveHackerNewsItem(ctx context.Context, id int) (*HackerNewsResponse, error) {
	return HackerNewsItem(id)
}

type ElasticsearchActivities struct{}

func (i *ElasticsearchActivities) Index(ctx context.Context, document *ElasticSearchDocument) error {
	return IndexDocument(*document)
}

type SearchRequest struct {
	Query string
	Size  int
}

func (i *ElasticsearchActivities) Search(ctx context.Context, req *SearchRequest) ([]ElasticSearchDocument, error) {
	return SearchDocuments(req.Query, req.Size)
}
