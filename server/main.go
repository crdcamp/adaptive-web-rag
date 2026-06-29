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

var AppConfig *Config

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file from root directory")
	}

	AppConfig, err = LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//llamaClient := CreateLlamaClient(AppConfig.LlamaBaseURL+"/v1", AppConfig.LlamaAPIKey)
	weaviateClient := CreateWeaviateClient(AppConfig.WeaviateBaseURL)

	// GetCollection(weaviateClient, "philosophyCollection")

	// // Might add a parameter for a custom output location
	// GenerateSearchQuery(llamaClient, AppConfig.ChatModel, "Tell me about some philosophies involving existential dread")
	// UnloadModel(AppConfig.ChatModel)
	//DeleteCollectionRw(weaviateClient, "philosophyCollection")
	CreateCollectionRw(weaviateClient, "philosophyCollection", "A test collection containing information on existential dread.")
	// CallCrawlScript()
	splitCrawlResults := SplitCrawlResults("crawl_data/crawl_results.json")
	EmbedText(weaviateClient, "philosophyCollection", splitCrawlResults)
	UnloadModel(AppConfig.EmbedModel)
	//func SplitEmbedAndUploadText(){}
}

type Config struct {
	ChatModel          string
	EmbedModel         string
	LlamaBaseURL       string
	LlamaAPIKey        string
	WeaviateBaseURL    string
	WeaviateBaseURLAlt string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ChatModel:          os.Getenv("CHAT_MODEL"),
		EmbedModel:         os.Getenv("EMBED_MODEL"),
		LlamaBaseURL:       os.Getenv("LLAMA_BASE_URL"),
		LlamaAPIKey:        os.Getenv("LLAMA_API_KEY"),
		WeaviateBaseURL:    os.Getenv("WEAVIATE_BASE_URL"),
		WeaviateBaseURLAlt: os.Getenv("WEAVIATE_BASE_URL_ALT"),
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
