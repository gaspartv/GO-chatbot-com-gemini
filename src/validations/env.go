package validations

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	API_URL                 string
	Port                    string
	OpenAI_API_Key          string
	OpenAI_API_URL          string
	OpenAI_API_Model        string
	OpenAI_API_Voice        string
	OpenAI_API_Instructions string
	GenAI_API_Key           string
	GenAI_API_Model         string
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, prosseguindo com variáveis do ambiente...")
	}

	config := &Config{
		API_URL:                 getEnvOrPanic("API_URL"),
		Port:                    getEnvOrPanic("PORT"),
		OpenAI_API_Key:          getEnvOrPanic("OPENAI_API_KEY"),
		OpenAI_API_URL:          getEnvOrPanic("OPENAI_API_URL"),
		OpenAI_API_Model:        getEnvOrPanic("OPENAI_API_MODEL"),
		OpenAI_API_Voice:        getEnvOrPanic("OPENAI_API_VOICE"),
		OpenAI_API_Instructions: getEnvOrPanic("OPENAI_API_INSTRUCTIONS"),
		GenAI_API_Key:           getEnvOrPanic("GENAI_API_KEY"),
		GenAI_API_Model:         getEnvOrPanic("GENAI_API_MODEL"),
	}

	return config
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Variável de ambiente obrigatória não definida: %s", key)
	}
	return value
}
