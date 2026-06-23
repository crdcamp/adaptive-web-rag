package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// 1. Declare the variable globally, but DO NOT initialize it here.
var AppConfig *Config

func main() {
	// 2. Load the .env file.
	// If you run 'go run main.go' from INSIDE the server folder, use "../.env"
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file from root directory")
	}

	// 3. Initialize the configuration NOW that the environment variables exist!
	AppConfig, err = LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println("APP CONFIG:", AppConfig.LlamaBaseURL)
	llamaClient := CreateLlamaClient(AppConfig.LlamaBaseURL+"/v1", AppConfig.LlamaAPIKey)
	weaviateClient := CreateWeaviateClient(AppConfig.WeaviateBaseURL)

	GetCollection(weaviateClient, "philosophyCollection")
	GenerateSearchQuery(llamaClient, AppConfig.ChatModel, "Tell me about some philosophies involving existential dread")
	UnloadModel(AppConfig.ChatModel)
	SplitCrawlResults("crawl_data/crawl_results.json")
}

type Config struct {
	ChatModel       string
	EmbedModel      string
	LlamaBaseURL    string
	WeaviateBaseURL string
	LlamaAPIKey     string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ChatModel:       os.Getenv("CHAT_MODEL"),
		EmbedModel:      os.Getenv("EMBED_MODEL"),
		LlamaBaseURL:    os.Getenv("LLAMA_BASE_URL"),
		WeaviateBaseURL: os.Getenv("WEAVIATE_BASE_URL"),
		LlamaAPIKey:     os.Getenv("LLAMA_API_KEY"),
	}

	return cfg, nil
}

// Create and return an OpenAI API compatible client for llama-server.
func CreateLlamaClient(baseURL string, apiKey string) openai.Client {
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		// API Key is not required for llama-server
		option.WithAPIKey(apiKey),
	)

	return client
}

// Create and return a Weaviate client for your Weaviate vector database.
func CreateWeaviateClient(host string) *weaviate.Client {
	cfg := weaviate.Config{
		Host:    host,
		Scheme:  "http",
		Headers: nil,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// Check the connection
	live, err := client.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Weaviate client live? %v\n", live)

	return client
}
