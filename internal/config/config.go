package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	TaskQueueName          string
	OllamaBaseURL          string
	OllamaModel            string
	GeminiApiKey           string
	GeminiModel            string
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
			OllamaBaseURL: getEnv("OLLAMA_BASE_URL", "http://localhost:11434/v1"),
			OllamaModel:   getEnv("OLLAMA_MODEL", "gemma3:latest"),
			GeminiApiKey:  getEnv("GEMINI_API_KEY", ""),
			GeminiModel:   getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
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
