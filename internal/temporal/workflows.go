package temporal

import (
	"fmt"
	"time"

	"github.com/shayansm2/temporallm/internal/chatbot"
	"github.com/shayansm2/temporallm/internal/elasticsearch"
	"go.temporal.io/sdk/workflow"
)

func SimpleChatWorkflow(ctx workflow.Context, message string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var llmActivities *LLMActivities

	var response string
	err := workflow.ExecuteActivity(ctx, llmActivities.SimpleChat, message).Get(ctx, &response)
	if err != nil {
		return "", fmt.Errorf("failed to get response from llm: %s", err)
	}
	return response, nil
}

func IndexHackerNewsStoryWorkflow(ctx workflow.Context, id int) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var hnActivities *HackerNewsApiActivities
	var esActivities *ElasticsearchActivities

	var storyTitle string
	items := []int{id}
	for len(items) > 0 {
		item := items[0]
		items = items[1:]
		var hnResponse *chatbot.HackerNewsResponse
		err := workflow.ExecuteActivity(ctx, hnActivities.RetrieveHackerNewsItem, item).Get(ctx, &hnResponse)
		if err != nil {
			return fmt.Errorf("failed to retrieve hacker news for item %d: %s", item, err)
		}
		if hnResponse.Title != "" {
			storyTitle = hnResponse.Title
		}
		doc := &elasticsearch.ESDocument{
			Id:    hnResponse.Id,
			Score: hnResponse.Score,
			Title: storyTitle,
			Type:  hnResponse.Type,
			Text:  hnResponse.Text,
		}
		err = workflow.ExecuteActivity(ctx, esActivities.Index, doc).Get(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to index hacker news item %d: %s", item, err)
		}
		items = append(items, hnResponse.Kids...)
	}
	return nil
}

func RetrivalAugmentedGenerationWorkflow(ctx workflow.Context, message string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var esActivities *ElasticsearchActivities
	var esDocs []elasticsearch.ESDocument
	err := workflow.ExecuteActivity(ctx, esActivities.Search, &SearchRequest{Query: message, Size: 5}).Get(ctx, &esDocs)
	if err != nil {
		return "", fmt.Errorf("failed to search on elasticsearch: %s", err)
	}

	var llmActivities *LLMActivities
	var response string
	systemMsg := chatbot.BuildSystemPrompt(esDocs)
	err = workflow.ExecuteActivity(ctx, llmActivities.Chat, message, systemMsg).Get(ctx, &response)
	if err != nil {
		return "", fmt.Errorf("failed to get response from llm: %s", err)
	}
	return response, nil
}
