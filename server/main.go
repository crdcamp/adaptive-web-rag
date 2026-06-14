package main

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// Dummy values for testing
const WeaviateEmbedURL string = "WEAVIATE-INCORRECT"
const ServerBaseURL string = "LLAMA-INCORRECT"

// These values need to be stored in a .env file I think (problem for later)
// Model Names
const ChatModel string = "Qwen2.5-7B-Instruct-Q4_K_M"
const EmbedModel string = "Qwen3-Embedding-8B-Q5_K_M"

const GeminiModel string= "gemma-4-12B-it-Q8_0-MTP"

// URLs
const LlamaBaseUrl string = "http://127.0.0.1:8080"
const WeaviateBaseUrl string = "http://127.0.0.1:8081"

const APIKey string = "no-key"

// We'll store these elsewhere later (probably a .env file... or .env file mixed with a toml or json config file)
// Just gotta get shit working for now so we're leaving it as is

// Might be better to directly pull these values from the `curl http://localhost:8001/v1/models | jq` results instead
// What does this mean? ^

// Need to address the error handling everywhere. We'll leave as is for now
func main() {
	weaviateClient := CreateWeaviateClient("localhost:8081")
	//DeleteCollection(weaviateClient, "CrawlResults")
	CreateCollection(weaviateClient, "CrawlResults", "A collection for storing internet results from web scraping")
	GenerateSearchQuery(ChatModel, "Tell me what it means to truly live")
}

func CreateWeaviateClient(host string) *weaviate.Client {
	// Weaviate client
	cfg := weaviate.Config{
		Host:    host,
		Scheme:  "http",
		Headers: nil,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}

	// Check the connection
	live, err := client.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Weaviate client live? %v\n", live)

	return client
}
