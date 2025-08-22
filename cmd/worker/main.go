package main

import (
	"log"

	"github.com/shayansm2/temporallm/internal/config"
	"github.com/shayansm2/temporallm/internal/llm"
	"github.com/shayansm2/temporallm/internal/temporal"
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

	w := worker.New(c, config.Load().TaskQueueName, worker.Options{})
	llm := llm.NewLLM(config.Load().OllamaBaseURL, "", config.Load().OllamaModel)
	w.RegisterWorkflow(temporal.SimpleChatWorkflow)
	w.RegisterWorkflow(temporal.IndexHackerNewsStoryWorkflow)
	w.RegisterWorkflow(temporal.RetrivalAugmentedGenerationWorkflow)
	w.RegisterActivity(&temporal.LLMActivities{LLM: &llm})
	w.RegisterActivity(&temporal.HackerNewsApiActivities{})
	w.RegisterActivity(&temporal.ElasticsearchActivities{})

	// Start the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
