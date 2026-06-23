package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

const ChatModel string = "Qwen2.5-7B-Instruct-Q4_K_M"
const EmbedModel string = "Qwen3-Embedding-8B-Q5_K_M"
const LlamaBaseUrl string = "http://127.0.0.1:8080"
const WeaviateBaseUrl string = "http://127.0.0.1:8081"
const APIKey string = "not-needed"

func main() {
	llamaClient := CreateLlamaClient(LlamaBaseUrl+"/v1", APIKey)
	//weaviateClient := CreateWeaviateClient("localhost:8081")

	// Function testing
	//DeleteCollection(weaviateClient, "testCollection")
	//CreateCollectionRw(weaviateClient, "philosophyCollection", "A collection for storing internet results from web scraping relating to philosophies on existential dread")
	GenerateSearchQuery(llamaClient, ChatModel, "Tell me about some philosophies involving existential dread")
	UnloadModel(ChatModel)
	//CallCrawlScript()
	SplitCrawlResults("crawl_data/crawl_results.json")
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
