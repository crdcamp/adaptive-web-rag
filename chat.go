package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

// A struct for handling HTTP requests for the server's chat interface.
// I will expand this to include time of request and maybe some other things
// in a later project.
type ChatPost struct {
	UserPrompt string
}

// These functions are looking a bit inefficient (I don't have exact reasoning for that at
// the moment). I think UserPrompt could just be declared once, but we'll look into that
// when the rough drafts of the functions are complete.

// Provide an introductory response to a user's prompt and determine the
// validity of their request. If the model deems the prompt to not be
// research related, the model respond by semantically redirecting the user.
// If the model deems the prompt to be research related, it will notify the user
// that it's checking the vector database for relevant results (but not
// actually check the vector database).
func InitialResponse(w http.ResponseWriter, chatPost ChatPost) string {
	userPrompt := chatPost.UserPrompt
	WriteBytes(w, "Introductory user prompt: "+userPrompt+"\n")

	InitialResponseSysPrompt := ReadMDFile("prompts/InitialResponseSysPrompt.md")
	InitialResponse := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, InitialResponseSysPrompt, userPrompt)
	WriteBytes(w, "Initial prompt response: "+InitialResponse+"\n")

	return InitialResponse
}

// Evaluate a prompt to determine if it valid or invalid to the scope
// of a research question. The model will strictly respond with `VALID`
// or `INVALID`.
func InitialPromptValidation(w http.ResponseWriter, chatPost ChatPost) string {
	userPrompt := chatPost.UserPrompt

	InitialPromptValidationSysPrompt := ReadMDFile("prompts/InitialPromptValidationSysPrompt.md")
	InitialPromptValidation := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, InitialPromptValidationSysPrompt, userPrompt)
	WriteBytes(w, "Prompt validation result: "+InitialPromptValidation+"\n")

	return InitialPromptValidation
}

func ValidPromptHandling(w http.ResponseWriter, chatPost ChatPost) {
	userPrompt := chatPost.UserPrompt

	searchQuery := GenerateSearchQuery(LlamaClient, AppConfig.ChatModelNoThink, userPrompt)
	WriteBytes(w, "Search query generated: "+searchQuery+"\nSearching memory...\n")
	WriteBytes(w, "Search query generated: "+searchQuery+"\nSearching memory...\n")

	// WWEEEEEEEP WEEEEEEP WEEEEEP WOOOOOOP EMERGENCY BELOW!!!!!!
	// `NearTextSearch` shouldn't be returning an error. `NearTextSearch` should be handling this on its own... Common bruv
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
		// Huge risk of prompt injection here (not a genuine concern)
		// This entire approach in general is pretty half assed awktually.. Let's just get it working
		WriteBytes(w, "Analyzing relevance of content for source: "+r.Source+"\n")

		vectorResponseValidationSysPrompt := ReadMDFile("prompts/vectorResponseValidationSysPrompt.md")
		vectorResponseUserPrompt := "You are given the following prompt :\n\n" + userPrompt +
			"\n\nBased on the given prompt, strictly classify the following context as RELEVANT or IRRELEVANT:\n" + r.Content

		vectorResponseValidation := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, vectorResponseValidationSysPrompt, vectorResponseUserPrompt)

		// Needs to be finished
		// Probably gonna end up being a switch statement
		WriteBytes(w, "Content relevancy conclusion: "+vectorResponseValidation+"\n")
		if vectorResponseValidation == "RELEVANT" {
			resultsSlice = append(resultsSlice, "SOURCE: "+r.Source,
				"CONTENT: "+r.Content+
					"\nEND OF CONTENT\n")
		}
	}

	// If allVectorDBResults is empty, call a function for internet search
	// that originates in llama-server.go
	if len(resultsSlice) != 0 {
		allVectorDBResults := strings.Join(resultsSlice, "\n")
		WriteBytes(w, "ALL VECTOR DB RESULTS:\n"+allVectorDBResults+"\n")

		vectorDBAnswer := AnswerWithResults(LlamaClient, userPrompt, allVectorDBResults)
		WriteBytes(w, "\n\n\nANSWER:\n")
		WriteBytes(w, vectorDBAnswer)
	} // ELSE CALL THE CRAWL SCRIPT AND RERUN THE LOOP AND STORE RELEVANT RESULTS
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
		// THIS SHOULD BE A FATAL ERROR. THE PROGRAM SHOULD STOP IF THIS BASIC CHECK HAS FAILED!
		slog.Warn("Error evaluating prompt. The model did not return the expected result of `VALID` or `INVALID`", "Here's the model's output: ", InitialPromptValidation)
	}
}
