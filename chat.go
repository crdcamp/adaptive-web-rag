package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/openai/openai-go/v3"
)

// A struct for handling HTTP requests for the server's chat interface.
// I will expand this to include time of request and maybe some other things
// in a later project.
type ChatPost struct {
	UserPrompt string
}

// THE SEARCH QUERY IS BEING GENERATED TWICE!!!! Or the update for it is being printed twice...
// I can't seem to figure out why (might have something to do with the mysterious black box that is Docker... probably not though)
// ALSO the invalid prompt responses are being printed twice...
// No matter, this can be fixed but not right now.

// Provide an introductory response to a user's prompt and determine the
// validity of their request. If the model deems the prompt to not be
// research related, the model will respond by semantically redirecting the user.
// If the model deems the prompt to be research related, it will notify the user
// that it's checking the vector database for relevant results (but not
// actually check the vector database).
func InitialResponse(w http.ResponseWriter, chatPost ChatPost) string {
	userPrompt := chatPost.UserPrompt
	WriteBytes(w, "Introductory user prompt: "+userPrompt+"\n")

	initialResponseSysPrompt := ReadMDFile("prompts/initialPromptSysPrompt.md")
	initialResponse := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, initialResponseSysPrompt, userPrompt)
	WriteBytes(w, "Initial prompt response: "+initialResponse+"\n")

	return initialResponse
}

// Evaluate a prompt to determine if it is valid or invalid in the context of
// research related questions. The model should strictly output `VALID`
// or `INVALID`.
func InitialPromptValidation(w http.ResponseWriter, chatPost ChatPost) string {
	userPrompt := chatPost.UserPrompt

	initialPromptValidationSysPrompt := ReadMDFile("prompts/initialPromptValidationSysPrompt.md")
	initialPromptValidation := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, initialPromptValidationSysPrompt, userPrompt)
	WriteBytes(w, "Prompt validation result: "+initialPromptValidation+"\n")

	return initialPromptValidation
}

// Analyze the relevance of given content in relation to the user's prompt.
// Can be used for vector database and search result validation.
func AnalyzeContentRelevance(w http.ResponseWriter, llamaClient openai.Client, chatPost ChatPost, contentToAnalyze string) string {
	userPrompt := chatPost.UserPrompt
	WriteBytes(w, "Analyzing relevance of content: "+contentToAnalyze+"\nEND OF CONTENT\n\n")

	vectorResponseValidationSysPrompt := ReadMDFile("prompts/vectorResponseValidationSysPrompt.md")
	vectorResponseUserPrompt := "You are given the following prompt: " + userPrompt +
		"\n\nBased on the given prompt, strictly classify the following context as RELEVANT or IRRELEVANT:\n" + contentToAnalyze

	vectorResponseValidation := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, vectorResponseValidationSysPrompt, vectorResponseUserPrompt)
	WriteBytes(w, "Content relevancy conclusion: "+vectorResponseValidation+"\n\n")

	return vectorResponseValidation
}

// If a vector DB search fails to provide relevant results, conduct an
// internet search and store relevant results in the vector DB.
func AnswerWithCrawlResults(w http.ResponseWriter, chatPost ChatPost, searchQuery string, crawlResultsFilePath string) {
	CallCrawlScript()
	crawlResults := ReadCrawlResults(crawlResultsFilePath)
	// content = crawlresults (find out how to parse through the json with something like this)
	for _, r := range crawlResults {
		relevance := AnalyzeContentRelevance(w, LlamaClient, chatPost, r.Content)
		///////////////////////////
		// TERRIBLE ERROR HANDLING/
		///////////////////////////
		switch relevance {
		case "RELEVANT":
			answer := AnswerWithVectorDBResults(LlamaClient, chatPost.UserPrompt, r.Content)
			WriteBytes(w, answer)
		case "IRRELEVANT":
			WriteBytes(w, "Irrelevant result\n")
		default:
			log.Fatal("FATAL ERROR EVALUATING PROMPT")
		}
	}
}

func ValidPromptHandling(w http.ResponseWriter, chatPost ChatPost) {
	userPrompt := chatPost.UserPrompt

	searchQuery := GenerateSearchQuery(LlamaClient, AppConfig.ChatModelNoThink, userPrompt)
	WriteBytes(w, "Search query generated: "+searchQuery+"\nSearching memory...\n")

	// WWEEEEEEEP WEEEEEEP WEEEEEP WOOOOOOP EMERGENCY BELOW!!!!!!
	// `NearTextSearch` shouldn't be returning an error. `NearTextSearch` should be handling the error on its own... Common bruv
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
	nearTextResults := nearTextStruct.Get[WebSearchCollection]
	var resultsSlice []string
	for _, r := range nearTextResults {
		vectorResponseValidation := AnalyzeContentRelevance(w, LlamaClient, chatPost, r.Content)
		if vectorResponseValidation == "RELEVANT" {
			resultsSlice = append(resultsSlice, "SOURCE: "+r.Source,
				"CONTENT: "+r.Content+"\nEND OF CONTENT\n")
		}
	}

	if len(resultsSlice) != 0 {
		allVectorDBResults := strings.Join(resultsSlice, "\n")
		vectorDBAnswer := AnswerWithVectorDBResults(LlamaClient, userPrompt, allVectorDBResults)
		WriteBytes(w, "\n\n\nANSWER:\n")
		WriteBytes(w, vectorDBAnswer)
	} else {
		AnswerWithCrawlResults(w, chatPost, searchQuery, "crawl_data/crawl_results.json")

		/////////////////////////////////////////////////
		// THEN ANALYZE THE RELEVANCE OF THE RESULTS AND
		// STORE ONLY THE RELEVANT ONES
		/////////////////////////////////////////////////
	}
}

func InvalidPromptHandling(w http.ResponseWriter, InitialResponse string, InitialPromptValidation string) {
	WriteBytes(w, "User prompt validation result: "+InitialPromptValidation+"\nInvalid prompt entered. Please ask a research-oriented question.\n")
	WriteBytes(w, "Model response: "+InitialResponse+"\n")
}

func InitialPromptDecisionTree(w http.ResponseWriter, chatPost ChatPost, InitialResponse string, InitialPromptValidation string) {
	switch InitialPromptValidation {
	case "VALID":
		ValidPromptHandling(w, chatPost)

	case "INVALID":
		InvalidPromptHandling(w, InitialResponse, InitialPromptValidation)

	default:
		log.Fatal("FATAL: Error evaluating prompt. The model did not return the expected result of `VALID` or `INVALID`", "Here's the model's output: ", InitialPromptValidation)
	}
}
