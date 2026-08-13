package main

import (
	"encoding/json"
	"io"
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

// Provide an introductory response to a user's prompt and determine the
// validity of their request. If the model deems the prompt to not be
// research related, the prompt will be semantically redirected. If the
// model deems the prompt to be research related, it will notify the user
// that it's checking the vector database for relevant results (but not
// actually check the vector database).
func initialResponse(chatPost ChatPost) string {
	userPrompt := chatPost.UserPrompt
	initialResponseSysPrompt := ReadMDFile("prompts/initialResponseSysPrompt.md")
	initialResponse := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, initialResponseSysPrompt, userPrompt)

	return initialResponse
}

// Evaluate a prompt to determine if it valid or invalid to the scope
// of a research question. The model will respond strictly with `VALID`
// or `INVALID`.
func initialPromptValidation(chatPost ChatPost) string {
	userPrompt := chatPost.UserPrompt
	initialResponseValidationSysPrompt := ReadMDFile("prompts/initialResponseValidationSysPrompt.md")
	initialResponseValidation := CreateChatCompletion(LlamaClient, AppConfig.ChatModelNoThink, initialResponseValidationSysPrompt, userPrompt)

	return initialResponseValidation
}

func initialPromptDecisionTree(initialResponseValidation string) string {
	return initialResponseValidation
}

// THIS NEEDS TO BE DIVIDED INTO SEVERAL FUNCTIONS
// Start by just verbally coding it out. Think of it like a decision tree my dude
func HandleChat(w http.ResponseWriter, r *http.Request) {
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

	initialResponse := initialResponse(chatPrompt)
	initialResponseValidation := initialResponseValidation(chatPrompt)

	// Initial response validation
	// THIS SHOULD PROBABLY BE ITS OWN FUNCTION
	switch initialResponseValidation {
	case "VALID":
		//Initial response
		searchQuery := GenerateSearchQuery(LlamaClient, AppConfig.ChatModelNoThink, chatPrompt.UserPrompt)
		WriteBytes(w, "User prompt validation result: "+initialResponseValidation+"\nSearch query: "+searchQuery+"\n")
		WriteBytes(w, "Model response: "+initialResponse+"\n")
		WriteBytes(w, "Searching memory...\n")

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
			vectorResponseUserPrompt := "You are given the following prompt :\n\n" + chatPrompt.UserPrompt +
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

			vectorDBAnswer := AnswerWithResults(LlamaClient, chatPrompt.UserPrompt, allVectorDBResults)
			WriteBytes(w, "\n\n\nANSWER:\n")
			WriteBytes(w, vectorDBAnswer)
		} // ELSE CALL THE CRAWL SCRIPT AND RERUN THE LOOP AND STORE RELEVANT RESULTS

	case "INVALID":
		//Initial response
		WriteBytes(w, "Model response: "+initialResponse+"\n")
		WriteBytes(w, "User prompt validation result: "+initialResponseValidation+"\nInvalid prompt entered. Please ask a research-oriented question.\n")

	// THIS SHOULD BE A FATAL ERROR. THE PROGRAM SHOULD STOP IF THIS BASIC CHECK HAS FAILED!
	default:
		slog.Warn("Error evaluating prompt. The model did not return the expected result of `VALID` or `INVALID`", "Here's the model's output: ", initialResponseValidation)
	}
}
