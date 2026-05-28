package main

// OpenAI API reference: https://pkg.go.dev/github.com/openai/openai-go/v3#section-readme

// INSTEAD OF SAVING THE JSON FILE,
// POST IT WITH HTTP AND THEN ACCESS THAT FROM THE
// PYTHON FILE

import (
	"bytes"
	"context" // A Context carries a deadline, a cancellation signal, and other values across API boundaries.
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
)

// {
//   "model": "ggml-org/gemma-3-4b-it-GGUF:Q4_K_M",
// }

// Unload a model from `llama-server`'s memory by sending a post request to the `/models/unload` endpoint.
func UnloadModel(modelName string) {
	const unloadURL = ServerBaseURL + "/models/unload"

	payload, _ := json.Marshal(map[string]string{"model": modelName})
	resp, err := http.Post(unloadURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err) // Probably need a different way to handle this error
	}
	defer resp.Body.Close()
	fmt.Printf("Status: %s\n", resp.Status)
}

// This occasionally outputs Mandarin characters...Luckily there's a solution for that but not a huge priority at the moment
// In the end, I'll probably just start personally converting models to GGUF format

// Also... alter this so it's just one search query for now. You can instead just pull more links from a single search instead of creating multiple searches
// The added complexity of the current method is unecessary and completely overkill
// In fact, you might consider making the vector database creation temporary (depending on how long embedding takes, that is)

// Also, there's some AI slop in this function as well but much of it follows the openai Go documentation
// But you know what that means... this needs some serious review
// When reviewing, refer to this specifically: https://pkg.go.dev/github.com/openai/openai-go/v3#readme-chat-completions-api
func GenerateSearchQueries(userPrompt string) {
	// Refer to this video for `context`: https://www.youtube.com/watch?v=BkzgYfygDy8
	ctx := context.Background() // Background returns a non-nil, empty Context. It is never canceled, has no values, and has no deadline. In other words, double check how to properly use this
	chatClient := openai.NewClient(
		option.WithBaseURL(ServerBaseURL),
		option.WithAPIKey(APIKey),
	)

	// I don't think it's necessary for this to be a variable dude
	systemPrompt := `You are a search query generator. When given a question or topic, generate exactly five search engine queries a person could enter into a browser to research it.

Respond ONLY with a JSON object in this exact format, no other text:
{"queries": ["query1", "query2", "query3", "query4", "query5"]}`

	resp, err := chatClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: ChatModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt), // Look into the difference between `DeveloperMessage()` and `SystemMessage()`
			openai.UserMessage(userPrompt),
		},
	})
	if err != nil {
		panic(err) // Need to double check if `panic` is a good idea here (almost definitely not)
	}

	raw := resp.Choices[0].Message.Content

	// More AI slop (the chat output format part). Will need to check this as well when I have more Go knowledge
	var result struct {
		Queries []string `json:"queries"`
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		panic(err)
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	println(string(out))   // Add fmt.Println() here
	UnloadModel(ChatModel) // Revisit UnloadModel. You can probably do all error handling within that function
}

//func callCrawlScript()
