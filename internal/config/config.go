package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	TaskQueueName          string
	OllamaBaseURL          string
	GeminiApiKey           string
	OpenAIApiKey           string
	LLMModel               string
	LLM                    string
	ElasticsearchURL       string
	ElasticsearchUser      string
	ElasticsearchPass      string
	ElasticsearchIndexName string
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		_ = godotenv.Overload()

		cfg = &Config{
			// temporal configs
			TaskQueueName: getEnv("TEMPORALTASK_QUEUE_NAME", "default"),
			// LLM configs
			LLM:           getEnv("LLM", "ollama"),
			OllamaBaseURL: getEnv("OLLAMA_BASE_URL", "http://localhost:11434/v1"),
			GeminiApiKey:  getEnv("GEMINI_API_KEY", ""),
			OpenAIApiKey:  getEnv("OPEN_AI_API_KEY", ""),
			LLMModel:      getEnv("LLM_MODEL", "gemma3:latest"),
			// knowledge base configs
			ElasticsearchURL:       getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
			ElasticsearchUser:      getEnv("ELASTICSEARCH_USER", "user"),
			ElasticsearchPass:      getEnv("ELASTICSEARCH_PASS", "pass"),
			ElasticsearchIndexName: getEnv("ELASTICSEARCH_INDEX_NAME", "hacker_news"),
		}
	})
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
