package temporal

import (
	"context"

	"github.com/shayansm2/temporallm/internal/chatbot"
	"github.com/shayansm2/temporallm/internal/elasticsearch"
	"github.com/shayansm2/temporallm/internal/llm"
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

func (i *HackerNewsApiActivities) RetrieveHackerNewsItem(ctx context.Context, id int) (*chatbot.HackerNewsResponse, error) {
	return chatbot.HackerNewsItem(id)
}

type ElasticsearchActivities struct{}

func (i *ElasticsearchActivities) Index(ctx context.Context, document *elasticsearch.ESDocument) error {
	return elasticsearch.IndexDocument(*document)
}

type SearchRequest struct {
	Query string
	Size  int
}

func (i *ElasticsearchActivities) Search(ctx context.Context, req *SearchRequest) ([]elasticsearch.ESDocument, error) {
	return elasticsearch.TextSearch(req.Query, req.Size)
}

func (i *ElasticsearchActivities) GetAllStories(ctx context.Context) ([]string, error) {
	return elasticsearch.GetAllStories()
}

func (i *ElasticsearchActivities) GetAllCommentsOfStory(ctx context.Context, title string) ([]string, error) {
	return elasticsearch.GetAllComments(title)
}
