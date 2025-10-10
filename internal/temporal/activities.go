package temporal

import (
	"context"

	"github.com/shayansm2/askhn/internal/elasticsearch"
	"github.com/shayansm2/askhn/internal/hackernews"
	"github.com/shayansm2/askhn/internal/llm"
)

type LLMActivities struct {
	LLM llm.LLM
}

func (i *LLMActivities) SimpleChat(ctx context.Context, message string) (string, error) {
	return i.LLM.Chat(llm.Message{User: message})
}

func (i *LLMActivities) Chat(ctx context.Context, userMsg, systemMsg string) (string, error) {
	return i.LLM.Chat(llm.Message{User: userMsg, System: systemMsg})
}

func (i *LLMActivities) AgenticChat(ctx context.Context, userMsg, systemMsg string) (llm.AgentResponse, error) {
	var response llm.AgentResponse
	err := llm.ChatWithSchema(i.LLM, llm.Message{User: userMsg, System: systemMsg}, &response)
	return response, err
}

type HackerNewsApiActivities struct{}

func (i *HackerNewsApiActivities) RetrieveHackerNewsItem(ctx context.Context, id int) (*hackernews.Item, error) {
	return hackernews.GetItem(id)
}

func (i *HackerNewsApiActivities) SearchHackerNews(ctx context.Context, query string) ([]int, error) {
	return hackernews.Search(query)
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
