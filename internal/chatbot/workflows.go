package chatbot

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SimpleChat(ctx workflow.Context, message string) (string, error) {
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

func IndexHackerNewsStory(ctx workflow.Context, id int) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var hnActivities *HackerNewsApiActivities
	var storyTitle string

	items := []int{id}
	for len(items) > 0 {
		item := items[0]
		items = items[1:]
		var hnResponse *HackerNewsResponse
		err := workflow.ExecuteActivity(ctx, hnActivities.RetrieveHackerNewsItem, item).Get(ctx, &hnResponse)
		if err != nil {
			return fmt.Errorf("failed to retrieve hacker news for item %d: %s", item, err)
		}
		if hnResponse.Title != "" {
			storyTitle = hnResponse.Title
		}
		req := IndexRequest{
			IndexName: ElasticsearchIndexName,
			Document: &ElasticSearchDocument{
				Id:    hnResponse.Id,
				Score: hnResponse.Score,
				Title: storyTitle,
				Type:  hnResponse.Type,
				Text:  hnResponse.Text,
			},
		}
		err = workflow.ExecuteActivity(ctx, hnActivities.IndexInElasticsearch, req).Get(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to index hacker news item %d: %s", item, err)
		}
		items = append(items, hnResponse.Kids...)
	}
	return nil
}
