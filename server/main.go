package main

import (
	"context"
	"fmt"
	"log"
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
	err := godotenv.Load("/.env")
	if err != nil {
		log.Fatal("Error loading .env file from root directory")
	}

	AppConfig, err = LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", DoHealthCheck).Methods("GET")
	log.Fatal(http.ListenAndServe(":8082", router))

}

func DoHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, i'm a golang microservice")
	w.WriteHeader(http.StatusAccepted) //RETURN HTTP CODE 202
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
