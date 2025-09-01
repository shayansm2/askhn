package temporal

import (
	"fmt"
	"strings"
	"time"

	"github.com/shayansm2/askhn/internal/chatbot"
	"github.com/shayansm2/askhn/internal/elasticsearch"
	"github.com/shayansm2/askhn/internal/llm"
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
	var step int
	workflow.SetQueryHandler(ctx, "steps", func() (map[string]interface{}, error) {
		return map[string]interface{}{
			"step":  step,
			"steps": []string{"kb_search", "llm"},
		}, nil
	})
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
	step++
	var llmActivities *LLMActivities
	var response string
	systemMsg := llm.SystemPromptBuilder{}.ForRAG(esDocs)
	err = workflow.ExecuteActivity(ctx, llmActivities.Chat, message, systemMsg).Get(ctx, &response)
	if err != nil {
		return "", fmt.Errorf("failed to get response from llm: %s", err)
	}
	step++
	return response, nil
}

type ProsConsRagParams struct {
	Message string
	Side    string
}

func ProsConsRagWorkflow(ctx workflow.Context, params ProsConsRagParams) (string, error) {
	var step int
	workflow.SetQueryHandler(ctx, "steps", func() (map[string]interface{}, error) {
		return map[string]interface{}{
			"step":  step,
			"steps": []string{"searching_kowledge_base", "generating_response_from_llm"},
		}, nil
	})
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var esActivities *ElasticsearchActivities
	var esDocs []elasticsearch.ESDocument
	err := workflow.ExecuteActivity(ctx, esActivities.Search, &SearchRequest{Query: params.Message, Size: 10}).Get(ctx, &esDocs)
	if err != nil {
		return "", fmt.Errorf("failed to search on elasticsearch: %s", err)
	}
	step++
	var contextBuilder strings.Builder
	for _, doc := range esDocs {
		contextBuilder.WriteString(fmt.Sprintf("title: %s\ncomment: %s\n\n", doc.Title, doc.Text))
	}
	context := strings.TrimSpace(contextBuilder.String())
	var llmActivities *LLMActivities
	var response string
	systemMsg, err := llm.GenerateSysPrompt(params.Side, map[string]string{
		"Context": context,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate system prompt: %s", err)
	}
	err = workflow.ExecuteActivity(ctx, llmActivities.Chat, params.Message, systemMsg).Get(ctx, &response)
	if err != nil {
		return "", fmt.Errorf("failed to get response from llm: %s", err)
	}
	step++
	return response, nil
}

func CreateGrandTruthDataWorkflow(ctx workflow.Context, path string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var esActivities *ElasticsearchActivities
	var stories []string
	err := workflow.ExecuteActivity(ctx, esActivities.GetAllStories).Get(ctx, &stories)
	if err != nil {
		return fmt.Errorf("failed to get all stories: %s", err)
	}
	for _, story := range stories {
		var comments []string
		err := workflow.ExecuteActivity(ctx, esActivities.GetAllCommentsOfStory, story).Get(ctx, &comments)
		if err != nil {
			return fmt.Errorf("failed to get all comments of story %s: %s", story, err)
		}
		fmt.Println(story, comments)
	}
	return nil
}
