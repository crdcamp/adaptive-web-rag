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
	GenerateSearchQuery(LlamaClient, AppConfig.ChatModelNoThink, TestQuery)
	CallCrawlScript()
	SplitEmbedAndUploadText(WeaviateClient, "crawl_data/crawl_results.json", WebSearchCollection)

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

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Welcome to the root page. Hit the `/chat` endpoint to chat, or hit the `/collections` endpoint to manage collections.\n"))
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}
}

//func writeBytes()

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
	initialResponseSysPrompt := `You are a research assistant with access to a vector database ("memory"). The vector database ("memory") contains results from past internet searches.

	When a user asks a question that could be answered using retrieved information, briefly note that you're checking memory before answering — one short line, not a ceremony (e.g. "Checking memory for prior notes on X."). Skip this step for trivial exchanges (greetings, clarifying questions, meta-questions about the conversation itself).

	DO NOT answer the user's prompt. You must only briefly not that you're checking your memory before answering.

	If a question is clearly outside your research scope (e.g. unrelated small talk, requests for content generation unrelated to the corpus, insults), say so plainly and redirect the user rather than attempting an answer anyway.`

	initialResponse := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, initialResponseSysPrompt, chatPrompt.UserPrompt)

	// Initial chat validation
	initialResponseValidationSysPrompt := `You are a classifier evaluating whether a given text is relevant to a specific research corpus.

	### Input
	<text_to_evaluate>
	[INSERT ANSWER OR PROMPT HERE]
	</text_to_evaluate>

	### Task
	Determine if the provided text is OUTSIDE the scope of research.

	OUT OF SCOPE includes:
	- Small talk, casual conversation, or insults (e.g., "How are you?", "Tell me a joke", "You're the worst")
	- Requests or content generation unrelated to the research subject
	- Insults, profanity, or conversational filler

	IN SCOPE includes:
	- Factual queries, analysis, or summaries relevant to the research corpus
	- Direct discussion of topics covered within the domain

	### Output Format
	Respond ONLY with either INVALID or VALID . Do not include any explanation, punctuation, or extra words.

	- Respond INVALID if the text is OUTSIDE the scope of research.
	- Respond VALID if the text is WITHIN the scope of research.`

	initialResponseValidation := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, initialResponseValidationSysPrompt, chatPrompt.UserPrompt)

	// Initial response validation
	switch initialResponseValidation {
	case "VALID":
		//Initial response
		searchQuery := GenerateSearchQuery(LlamaClient, AppConfig.ChatModelNoThink, chatPrompt.UserPrompt)
		_, err = w.Write([]byte("Response validation result: " + initialResponseValidation + "\nSearch query: " + searchQuery + "\n"))
		if err != nil {
			slog.Error("error writing response body", "err", err)
		}

		_, err = w.Write([]byte("Model response: " + initialResponse + "\n"))
		if err != nil {
			slog.Error("error writing response body", "err", err)
		}

		_, err = w.Write([]byte("Searching memory..."))
		if err != nil {
			slog.Error("error writing response body", "err", err)
		}

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
			if r.Content == "" {
				_, err := w.Write([]byte("NO CONTENT FOUND\n"))
				if err != nil {
					slog.Error("error writing response body", "err", err)
					return
				}
			} else {
				_, err := w.Write([]byte("CONTENT FOUND:\n" + r.Content + "\n\n"))
				if err != nil {
					slog.Error("error writing response body", "err", err)
					return
				}
			}
		}

	case "INVALID":
		//Initial response
		_, err = w.Write([]byte("Model response: " + initialResponse + "\n"))
		if err != nil {
			slog.Error("error writing response body", "err", err)
		}
		_, err = w.Write([]byte("Response validation result: " + initialResponseValidation + "\nInvalid prompt entered. Please ask a research-oriented question.\n"))
		if err != nil {
			slog.Error("error writing response body", "err", err)
		}

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
