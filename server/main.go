package main

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// A bunch of variables that I'll store correctly later
const ChatModel string = "Qwen2.5-7B-Instruct-Q4_K_M"
const EmbedModel string = "Qwen3-Embedding-8B-Q5_K_M"
const ServerBaseURL string = "http://127.0.0.1:8001"
const APIKey string = "no-key"

// We'll store these elsewhere later (probably a .env file... or .env file mixed with a toml or json config file)
// Just gotta get shit working for now so we're leaving it as is

// Might be better to directly pull these values from the `curl http://localhost:8001/v1/models | jq` results instead
// What does this mean? ^

// Need to address the error handling everywhere. We'll leave as is for now
func main() {
	// Weaviate client
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	weaviateClient, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// Weaviate status check
	live, err := weaviateClient.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", live)

	CreateCollection(weaviateClient, "CrawlResults", "A collection for storing internet results from web scraping")
	GenerateSearchQuery(ChatModel, "Tell me about the benefits and drawbacks of llama.cpp")
	//CallCrawlScript()
	//ChunkEmbedAndUploadCrawlResults()
}

// Can probably just remove this one
//const WeaviateClient
