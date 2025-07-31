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
