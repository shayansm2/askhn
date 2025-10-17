package main

import (
	"log"

	"github.com/shayansm2/askhn/internal/config"
	"github.com/shayansm2/askhn/internal/elasticsearch"
	"github.com/shayansm2/askhn/internal/llm"
	"github.com/shayansm2/askhn/internal/temporal"
	"go.temporal.io/sdk/worker"
)

func getLLM(cnf *config.Config) llm.LLM {
	switch cnf.LLM {
	case "openai":
		llm := llm.NewOpenAiLLM(cnf.OpenAIApiKey, cnf.LLMModel)
		return &llm
	case "gemini":
		llm := llm.NewGeminiLLM(cnf.GeminiApiKey, cnf.LLMModel)
		return &llm
	case "ollama":
		llm := llm.NewOllamaLLM(cnf.OllamaBaseURL+"/v1", cnf.LLMModel)
		return &llm
	default:
		panic("LLM env variable should be one of openai, gemini or ollama")
	}
}

func main() {
	cnf := config.Load()
	c := temporal.GetClient()

	w := worker.New(c, cnf.TaskQueueName, worker.Options{})
	llm := getLLM(cnf)
	elasticsearch.CreateElasticsearchIndex()

	w.RegisterWorkflow(temporal.SimpleChatWorkflow)
	w.RegisterWorkflow(temporal.IndexHackerNewsStoryWorkflow)
	w.RegisterWorkflow(temporal.RetrivalAugmentedGenerationWorkflow)
	w.RegisterWorkflow(temporal.ProsConsRagWorkflow)
	w.RegisterWorkflow(temporal.AgenticRAGWorkflow)

	w.RegisterActivity(&temporal.LLMActivities{LLM: llm})
	w.RegisterActivity(&temporal.HackerNewsApiActivities{})
	w.RegisterActivity(&temporal.ElasticsearchActivities{})

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
