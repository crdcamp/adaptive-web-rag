package main

import (
	"bytes"
	"context"
	"encoding/json"
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

	llamaClient := CreateLlamaClient(AppConfig.LlamaBaseURL, AppConfig.LlamaAPIKey)

	mux := mux.NewRouter()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/chat/", handleChat)
	log.Fatal(http.ListenAndServe(":8082", mux))
}

type ChatPost struct {
	UserPrompt string
}

func handleRoot(w http.ResponseWriter, _ *http.Request, llamaClient ) {
	_, err := w.Write([]byte("Welcome to the root page. Hit the chat endpoint instead please.\n"))
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}
}

func handleChat(w http.ResponseWriter, , llamaClient openai.Client, chatPrompt string) {
	systemPrompt := "You are a helpful assistant. Answer the question to the best of your ability"
	chatResponse := CreateChatCompletion(llamaClient, AppConfig.ChatModel, systemPrompt, chatPrompt)

	var output bytes.Buffer
	output.WriteString(chatResponse)

	_, err := w.Write(output.Bytes())
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}

func handleJSON(w http.ResponseWriter, r *http.Request, llamaClient openai.Client) {
	ctx := context.Background() // Add a timeout to this

	byteData, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading request body", "err", err)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	var chatPrompt ChatPost
	err := json.Unmarshal(byteData, &chatPrompt)
	if err != nil {
		slog.Error("error unmarshalling chat request body", "err", err)
		http.Error(w, "error parsing request JSON", http.StatusBadRequest)
		return
	}

	if reqData.Name == "" {
		http.Error(w, "no input provided", http.StatusBadRequest)
		return
	}
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
