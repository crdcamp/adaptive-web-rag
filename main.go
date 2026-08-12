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
	mux.HandleFunc("/chat", handleChat)
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

// Will expand to include time of request and maybe some other things
type ChatPost struct {
	UserPrompt string
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

func writeBytes(w http.ResponseWriter, input string) {
	_, err := w.Write([]byte(input))
	if err != nil {
		slog.Error("error writing response body", "err", err)
	}
}

// Needs a file extension check
func readMDFile(filePath string) string {
	resultBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	resultString := string(resultBytes)
	return resultString
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	writeBytes(w, "Welcome to the root page. Hit the `/chat` endpoint to chat, or hit the `/collections` endpoint to manage collections.\n")
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

	// Initial  chat response
	initialResponseSysPromptData := readMDFile("prompts/initialResponseSysPrompt.md")
	initialResponseSysPrompt := initialResponseSysPromptData
	fmt.Printf("initialResponseSysPrompt result:\n%v", initialResponseSysPrompt)
	initialResponse := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, initialResponseSysPrompt, chatPrompt.UserPrompt)

	// Initial chat validation
	initialResponseValidationSysPromptData := readMDFile("prompts/initialResponseValidationSysPrompt.md")
	initialResponseValidationSysPrompt := initialResponseValidationSysPromptData
	initialResponseValidation := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, initialResponseValidationSysPrompt, chatPrompt.UserPrompt)

	// Initial response validation
	switch initialResponseValidation {
	case "VALID":
		//Initial response
		searchQuery := GenerateSearchQuery(LlamaClient, AppConfig.ChatModelNoThink, chatPrompt.UserPrompt)
		writeBytes(w, "Response validation result: "+initialResponseValidation+"\nSearch query: "+searchQuery+"\n")
		writeBytes(w, "Model response: "+initialResponse+"\n")
		writeBytes(w, "Searching memory...")

		// This shouldn't be returning an error. NearTextSearch should be handling this on its own
		nearTextBytes, err := NearTextSearch(WeaviateClient, WebSearchCollection, 3, searchQuery)
		if err != nil {
			slog.Error("error running near text search", "err", err)
			http.Error(w, "error searching memory", http.StatusInternalServerError)
			return
		}

		var nearTextStruct HrefAndContentDBResponse
		if err := json.Unmarshal(nearTextBytes, &nearTextStruct); err != nil {
			slog.Error("error unmarshalling near text search response", "err", err)
			http.Error(w, "error parsing memory search results", http.StatusInternalServerError)
			return
		}

		// All the near text data handling desperately needs to be cleaned up my god
		// ... To be fair the GraphQL stuff is just ridiculously difficult to figure out...
		// I need to learn more about Go.
		nearTextResults := nearTextStruct.Get[WebSearchCollection]
		for _, r := range nearTextResults {
			// I'm almost certain this check isn't sufficient
			if r.Content == "" {
				writeBytes(w, "NO CONTENT FOUND\n")
			} else {
				writeBytes(w, "CONTENT FOUND (source: "+r.Source+"):\n"+r.Content+"\n\n")
				// Huge risk of prompt injection here (not a genuine concern for the current scope, but worth noting for now)
				// Another time you can create a function that does this to avoid that issue
				// This entire approach in general is pretty half assed awktually
				vectorResponseValidationSysPrompt := `You are a strict relevance classifier for a RAG pipeline.

				Your task is to determine whether the retrieved text context is relevant to the user's prompt.

				Context is RELEVANT if it contains direct answers, partial facts, or necessary background information to help address the prompt. Otherwise, it is IRRELEVANT.

				User Prompt:
				"""
				` + chatPrompt.UserPrompt + `
				"""

				Retrieved Context:
				"""
				` + r.Content + `
				"""

				Respond with EXACTLY one word: "RELEVANT" or "IRRELEVANT". Do not include quotes, explanations, or any other text.`
				vectorResponseValidation := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, vectorResponseValidationSysPrompt, r.Content)

				// Needs to be finished
				if vectorResponseValidation == "RELEVANT" {
					writeBytes(w, "Model response: "+initialResponse+"\n")
				}
			}
		}
	case "INVALID":
		//Initial response
		writeBytes(w, "Model response: "+initialResponse+"\n")
		writeBytes(w, "Response validation result: "+initialResponseValidation+"\nInvalid prompt entered. Please ask a research-oriented question.\n")

	default:
		slog.Warn("Error evaluating prompt. The model did not return the expected result of `VALID` or `INVALID`", "Here's the model's output: ", initialResponseValidation)
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
