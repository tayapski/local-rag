package config

import (
	"sync"
	"os"
)

type Config struct {
	OllamaURL	string
	QdrantHost	string
	QdrantPort	string
}

var (
	instance *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func () {
		instance = &Config{
			OllamaURL: getEnvWithDefault("OLLAMA_URL", "http://localhost:11434"),
			QdrantHost: getEnvWithDefault("QDRANT_HOST", "localhost"),
			QdrantPort: getEnvWithDefault("QDRANT_PORT", "6334"),
		}
	})
	return instance
}

func getEnvWithDefault(key string, defaultVal string) string {
	if val, isAssigned := os.LookupEnv(key); isAssigned {
		return val
	}
	return defaultVal
}