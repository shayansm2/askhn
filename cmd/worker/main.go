package main

import (
	"log"

	"github.com/shayansm2/askhn/internal/config"
	"github.com/shayansm2/askhn/internal/llm"
	"github.com/shayansm2/askhn/internal/temporal"
	"go.temporal.io/sdk/worker"
)

func main() {
	c := temporal.GetClient()

	w := worker.New(c, config.Load().TaskQueueName, worker.Options{})
	llm := llm.NewOllamaLLM(config.Load().OllamaBaseURL, config.Load().OllamaModel)
	// llm := llm.NewGeminiLLM(config.Load().GeminiApiKey, config.Load().GeminiModel)

	w.RegisterWorkflow(temporal.SimpleChatWorkflow)
	w.RegisterWorkflow(temporal.IndexHackerNewsStoryWorkflow)
	w.RegisterWorkflow(temporal.RetrivalAugmentedGenerationWorkflow)
	w.RegisterWorkflow(temporal.ProsConsRagWorkflow)

	w.RegisterActivity(&temporal.LLMActivities{LLM: &llm})
	w.RegisterActivity(&temporal.HackerNewsApiActivities{})
	w.RegisterActivity(&temporal.ElasticsearchActivities{})

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
