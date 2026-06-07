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
const WeaviateEmbedURL string = "http://llama-server:8080"
const WeaviateClientHost string = "localhost:8080"
const APIKey string = "no-key"

// We'll store these elsewhere later (probably a .env file... or .env file mixed with a toml or json config file)
// Just gotta get shit working for now so we're leaving it as is

// Might be better to directly pull these values from the `curl http://localhost:8001/v1/models | jq` results instead
// What does this mean? ^

// Need to address the error handling everywhere. We'll leave as is for now
func main() {
	// Weaviate client
	cfg := weaviate.Config{
		Host:   WeaviateClientHost,
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
	fmt.Println("WeaviateClient live and well?", live)

	// We'll turn this all into a function within the main function at some point
	//TestSplit()
	//UnloadModel(EmbedModel)
	//fmt.Println(string(GetCollection(weaviateClient,  "CrawlResults")))
	//GenerateSearchQuery(ChatModel, "Tell me about the benefits and drawbacks of using llama.cpp")
	//NearTextSearch(weaviateClient, "CrawlResults", 1, "What are the various ways llama.cpp can be used?")
	//runtime.GC() // Frees memory... sorta
	//debug.FreeOSMemory()
	DeleteCollection(weaviateClient, "CrawlResults")
	CreateCollection(weaviateClient, "CrawlResults", "A collection for storing internet results from web scraping")
	SplitEmbedAndUploadCrawlResults(weaviateClient, "CrawlResults")
	// CallCrawlScript()
	// SplitEmbedAndUploadCrawlResults(EmbedModel)
}
