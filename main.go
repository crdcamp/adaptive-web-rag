package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// Configurations (they're a mess)
var AppConfig *Config
var LlamaClient openai.Client
var WeaviateClient *weaviate.Client
var WebSearchCollection string
var TestQuery string

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}

	AppConfig, err = loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Clients for the models (LlamaClient) and the vector database (WeaviateClient)
	LlamaClient = createLlamaClient(AppConfig.LlamaURL, AppConfig.LlamaAPIKey)
	WeaviateClient = createWeaviateClient(AppConfig.WeaviateURL)

	WebSearchCollection = "WebSearchCollection"
	TestQuery = "What are some of the greatest mysteries throughout ancient history?"

	CreateCollection(WeaviateClient, WebSearchCollection, "A collection of web search results conducted and stored by an LLM.")
	// GenerateSearchQuery(LlamaClient, AppConfig.ChatModelNoThink, TestQuery)
	// CallCrawlScript()
	// SplitEmbedAndUploadText(WeaviateClient, "crawl_data/crawl_results.json", WebSearchCollection)

	// Endpoint for the actual chat
	mux := mux.NewRouter()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/chat", HandleChat)
	// Future endpoint for collections management
	// mux.HandleFunc("/collections", handleCollections) Need this later

	log.Fatal(http.ListenAndServe(":8082", mux))
}

type Config struct {
	LlamaURL         string
	LlamaAPIKey      string
	WeaviateURL      string
	ChatModelThink   string
	ChatModelNoThink string
	EmbedModel       string
}

// Load the .env file variables.
func loadConfig() (*Config, error) {
	cfg := &Config{
		LlamaURL:         os.Getenv("LLAMA_URL"),
		LlamaAPIKey:      os.Getenv("LLAMA_API_KEY"),
		WeaviateURL:      os.Getenv("WEAVIATE_URL"),
		ChatModelThink:   os.Getenv("CHAT_MODEL_THINK"),
		ChatModelNoThink: os.Getenv("CHAT_MODEL_NO_THINK"),
		EmbedModel:       os.Getenv("EMBED_MODEL"),
	}

	return cfg, nil
}

func WriteBytes(w http.ResponseWriter, input string) {
	_, err := w.Write([]byte(input))
	if err != nil {
		slog.Error("error writing response body", "err", err)
	}
}

// Needs a file extension check
func ReadMDFile(filePath string) string {
	resultBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	resultString := string(resultBytes)
	return resultString
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	WriteBytes(w, "Welcome to the root page. Hit the `/chat` endpoint to chat, or hit the `/collections` endpoint to manage collections.\n(Note: Collections management is not implemented.)")
}

// Create and return an OpenAI API compatible client for llama-server.
func createLlamaClient(baseURL string, apiKey string) openai.Client {
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		// API Key is not required for llama-server (but hey, ya neva know)
		option.WithAPIKey(apiKey),
	)

	return client
}

// Create and return a Weaviate client for your Weaviate vector database.
func createWeaviateClient(host string) *weaviate.Client {
	cfg := weaviate.Config{
		Host:   host,
		Scheme: "http",
		Headers: map[string]string{
			"X-Openai-Api-Key": AppConfig.LlamaAPIKey,
		},
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
