package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

var AppConfig, _ = LoadConfig()

func main() {
	fmt.Println("APP CONFIG:", AppConfig.LlamaBaseURL)
	llamaClient := CreateLlamaClient(AppConfig.LlamaBaseURL+"/v1", AppConfig.LlamaAPIKey)
	//weaviateClient := CreateWeaviateClient("localhost:8081")
	//DeleteCollection(weaviateClient, "testCollection")
	//CreateCollectionRw(weaviateClient, "philosophyCollection", "A collection for storing internet results from web scraping relating to philosophies on existential dread")
	GenerateSearchQuery(llamaClient, AppConfig.ChatModel, "Tell me about some philosophies involving existential dread")
	UnloadModel(AppConfig.ChatModel)
	//CallCrawlScript()
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
