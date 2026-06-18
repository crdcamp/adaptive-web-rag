package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// Dummy values that we'll get rid of when rewrite branch is complete
const WeaviateEmbedURL string = "WEAVIATE-INCORRECT"
const ServerBaseURL string = "LLAMA-INCORRECT"

// We'll store these elsewhere later (probably a .env file... or .env file mixed with a toml or json config file)
// Just gotta get shit working for now so we're leaving it as is
const ChatModel string = "Qwen2.5-7B-Instruct-Q4_K_M"
const EmbedModel string = "Qwen3-Embedding-8B-Q5_K_M"
const LlamaBaseUrl string = "http://127.0.0.1:8080/v1"
const WeaviateBaseUrl string = "http://127.0.0.1:8081"
const APIKey string = "not-needed"

func main() {
	llamaClient := CreateLlamaClient("http://localhost:8080/v1", APIKey)
	weaviateClient := CreateWeaviateClient("localhost:8081")

	// Check if this stuff works
	CreateCollectionRw(weaviateClient, "testCollection", "A collection for storing internet results from web scraping")
	GenerateSearchQuery(llamaClient, ChatModel, "Tell me about some philosophies involving existential dread")
	UnloadModel(ChatModel)
	DeleteCollectionRw(weaviateClient, "testCollection")

	// Check if crawl script works
	CallCrawlScriptRw()
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

// Calls crawl.py to conduct web search. Results are saved to `server/crawl_data/crawl_results.json`.
func CallCrawlScriptRw() {
	// Run crawl.py
	fmt.Println("Executing: crawl.py")
	cmd := exec.Command("python3", "crawl.py")

	// Output to terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal("Error when running command: ", err)
	} else {
		fmt.Println("Successfully executed: crawl.py")
	}
}
