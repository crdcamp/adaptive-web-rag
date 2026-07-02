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
	"github.com/weaviate/weaviate/entities/models"
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

	// Testing
	var webSearchCollection string = "WebSearchCollection"
	// var webSearchCollectionDescription string = "A collection of various web searches produced by the model."
	var prompt string = "What are the projected API costs for LLMs in the next decade"

	// THIS WILL BE A FUNCTION
	//CreateCollection(weaviateClient, webSearchCollection, "A collection of various web searches produced by the model.")
	// GenerateSearchQuery(llamaClient, AppConfig.ChatModel, prompt)
	// CallCrawlScript()
	// splitCrawlResults := SplitCrawlResults("crawl_data/crawl_results.json")
	// EmbedText(weaviateClient, webSearchCollection, splitCrawlResults)

	// fmt.Println(ReadAllCollectionNames(weaviateClient))
	// NearTextSearch(weaviateClient, "PhilosophyCollection", 4, "If nothing I do matters in a million years, does it matter now?")
	//CreateChatCompletion(llamaClient, AppConfig.ChatModel, "Answer the question to the best of your abilities", "What are the best use cases for a vector database?")
	//GenerateSearchQuery(llamaClient, AppConfig.ChatModel, "Tell me about philosophies involving existential dread")
	// CreateCollection(weaviateClient, webSearchCollection, webSearchCollectionDescription)
	// GenerateSearchQuery(llamaClient, AppConfig.ChatModel, prompt)
	// CallCrawlScript()
	// splitCrawlResults := SplitCrawlResults("crawl_data/crawl_results.json")
	// EmbedText(weaviateClient, webSearchCollection, splitCrawlResults)

	// This is a mess
	vectorSearchResult := vectorSearch(*weaviateClient, llamaClient, webSearchCollection, prompt)
	vectorContent := vectorSearchResult.Data
	fmt.Println(vectorContent)
	//AnswerWithVectorDBResults(llamaClient, string(vectorSearchResult))
}

// I feel like there's a better way to do this than with nested functions
func internetSearch(llamaClient openai.Client, weaviateClient *weaviate.Client, prompt string) {
	GenerateSearchQuery(llamaClient, AppConfig.ChatModel, prompt)
	UnloadModel(AppConfig.ChatModel)
	// Should probably add an output location for this
	CallCrawlScript()
}

type Config struct {
	ChatModel          string
	EmbedModel         string
	LlamaBaseURL       string
	LlamaServer        string
	LlamaAPIKey        string
	WeaviateBaseURL    string
	WeaviateBaseURLAlt string
	// LlamaClient        *openai.Client
	// WeaviateClient     *weaviate.Client
}

func splitEmbedAndUploadText(weaviateClient *weaviate.Client, className string, crawlResultsPath string) {
	splitCrawlResults := SplitCrawlResults(crawlResultsPath)
	EmbedText(weaviateClient, "philosophyCollection", splitCrawlResults)
}

func vectorSearch(weaviateClient weaviate.Client, llamaClient openai.Client, className string, prompt string) *models.GraphQLResponse {
	fmt.Printf("Searching vector database for prompt: %q\n", prompt)
	query := RefineVectorSearchQuery(llamaClient, prompt)
	UnloadModel(AppConfig.ChatModel)

	return NearTextSearch(&weaviateClient, className, 3, query)
}

// Load the .env file variables.
func LoadConfig() (*Config, error) {
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
