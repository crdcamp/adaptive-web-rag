package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	// Endpoint for the actual chat
	mux := mux.NewRouter()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/chat", handleChat)
	// Future endpoint for collections management
	// mux.HandleFunc("/collections", handleCollections) Need this later

	log.Fatal(http.ListenAndServe(":8082", mux))
}

type Config struct {
	LlamaURL    string
	LlamaAPIKey string
	WeaviateURL string
	ChatModel   string
	EmbedModel  string
}

// Will expand to include time of request and maybe some other things
type ChatPost struct {
	UserPrompt string
}

// Load the .env file variables.
func loadConfig() (*Config, error) {
	cfg := &Config{
		LlamaURL:    os.Getenv("LLAMA_URL"),
		LlamaAPIKey: os.Getenv("LLAMA_API_KEY"),
		WeaviateURL: os.Getenv("WEAVIATE_URL"),
		ChatModel:   os.Getenv("CHAT_MODEL"),
		EmbedModel:  os.Getenv("EMBED_MODEL"),
	}

	return cfg, nil
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Welcome to the root page. Hit the `/chat` endpoint instead please.\n"))
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	// Intake the byte data provided by the request
	byteData, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading request body", "err", err)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	// Create a variable chatPrompt from ChatPost struct to unmarshal the byte data
	var chatPrompt ChatPost
	err = json.Unmarshal(byteData, &chatPrompt)
	if err != nil {
		slog.Error("error unmarshalling chat request body", "err", err)
		http.Error(w, "error parsing request JSON", http.StatusBadRequest)
		return
	}

	if chatPrompt.UserPrompt == "" {
		http.Error(w, "No prompt was provided. Please try again", http.StatusBadRequest)
		return
	}

	// First we want to tell the user that we're checking the "memory" to see if
	// their question can be answered with internet results

	// .... Buuuuuut, being able to update the user is another thing I'd
	// have to learn so we're going to just check the memory for now

	// We need to do a half assed conversion to search the vector database
	//vectorDBQuerySysPrompt := "You are a search query generator. When given a question or topic, generate ONE search engine query that a person could enter into a browser to research it."
	//vectorDBQuery := CreateChatCompletion(LlamaClient, AppConfig.ChatModel, vectorDBQuerySysPrompt, chatPrompt.UserPrompt)
	searchQuery := GenerateSearchQuery(LlamaClient, AppConfig.ChatModel, chatPrompt.UserPrompt)

	// test := NearTextSearch(WeaviateClient, "WebSearchCollection", 3, vectorDBQuery)
	// fmt.Println(test)

	// use http POST to update user

	// Then, if the answer is in memory, answer using that.

	// If not, conduct an internet search

	// If the results are relevant to the question, then answer using those
	// If they aren't, inform the user you do not have extended search capability yet
	// or something like that.

	_, err = w.Write([]byte(searchQuery))
	if err != nil {
		slog.Error("error writing response body", "err", err)
	}
}

// Create and return an OpenAI API compatible client for llama-server.
func createLlamaClient(baseURL string, apiKey string) openai.Client {
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		// API Key is not required for llama-server
		option.WithAPIKey(apiKey),
	)

	return client
}

// Create and return a Weaviate client for your Weaviate vector database.
func createWeaviateClient(host string) *weaviate.Client {
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
