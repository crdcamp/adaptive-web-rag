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

	// Might be better to just have these in every function
	llamaClient := CreateLlamaClient(AppConfig.LlamaServer+"/v1", AppConfig.LlamaAPIKey)
	weaviateClient := CreateWeaviateClient(AppConfig.WeaviateBaseURL)

	CreateCollection(weaviateClient, "humanDiscoveries", "A collection for testing the RAG pipeline containing information on human discoveries")

	internetSearch(llamaClient, weaviateClient, "What are some of the greatest discoveries humanity has made?")
	splitEmbedAndUploadText(weaviateClient, "humandDiscovervies", "crawl_data/crawl_results.json")
}

// GetCollection(weaviateClient, "philosophyCollection")
// DeleteCollectionRw(weaviateClient, "philosophyCollection")
// CreateCollectionRw(weaviateClient, "philosophyCollection", "A test collection containing information on existential dread.")
func internetSearch(llamaClient openai.Client, weaviateClient *weaviate.Client, question string) {
	GenerateSearchQuery(llamaClient, AppConfig.ChatModel, question)
	UnloadModel(AppConfig.ChatModel)
	// Should probably add an output location for this
	CallCrawlScript()
}

func splitEmbedAndUploadText(weaviateClient *weaviate.Client, className string, crawlResultsPath string) {
	splitCrawlResults := SplitCrawlResults(crawlResultsPath)
	EmbedText(weaviateClient, "philosophyCollection", splitCrawlResults)
}

type Config struct {
	ChatModel          string
	EmbedModel         string
	LlamaBaseURL       string
	LlamaServer        string
	LlamaAPIKey        string
	WeaviateBaseURL    string
	WeaviateBaseURLAlt string
}

// Load the .env file variables.
func LoadConfig() (*Config, error) {
	// This maybe could be a loop?
	cfg := &Config{
		ChatModel:          os.Getenv("CHAT_MODEL"),
		EmbedModel:         os.Getenv("EMBED_MODEL"),
		LlamaBaseURL:       os.Getenv("LLAMA_BASE_URL"),
		LlamaServer:        os.Getenv("LLAMA_SERVER"),
		LlamaAPIKey:        os.Getenv("LLAMA_API_KEY"),
		WeaviateBaseURL:    os.Getenv("WEAVIATE_BASE_URL"),
		WeaviateBaseURLAlt: os.Getenv("WEAVIATE_BASE_URL_ALT"), // Will delete eventually
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
