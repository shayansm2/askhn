package main

import (
	"log"

	"github.com/shayansm2/temporallm/internal/chatbot"
	"github.com/shayansm2/temporallm/internal/utils/llm"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, chatbot.TaskQueueName, worker.Options{})
	llm := llm.NewLLM(chatbot.OllamaBaseURL, "", chatbot.OllamaModelGemma3)
	activities := &chatbot.LLMActivities{LLM: &llm}

	w.RegisterWorkflow(chatbot.SimpleChat)
	w.RegisterActivity(activities)

	// Start the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
