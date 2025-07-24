package main

import (
	"log"

	"github.com/shayansm2/temporallm"
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

	w := worker.New(c, temporallm.TaskQueueName, worker.Options{})
	llm := temporallm.NewLLM(temporallm.OllamaBaseURL, "", temporallm.OllamaModelGemma3)
	activities := &temporallm.LLMActivities{LLM: &llm}

	w.RegisterWorkflow(temporallm.SimpleChat)
	w.RegisterActivity(activities)

	// Start the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
