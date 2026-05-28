package main

import (
	"bytes"
	"context" // A Context carries a deadline, a cancellation signal, and other values across API boundaries.
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
)

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

// Generate a search query to pass into `crawl.py`
func GenerateSearchQueries(modelName string, userPrompt string) string {
	ctx := context.Background()
	client := openai.NewClient(
		option.WithBaseURL(ServerBaseURL),
		option.WithAPIKey(APIKey), // Leaving this here just in case
	)
	systemMessage := "You are a search query generator. When given a question or topic, generate a search engine query that a person could enter into a browser to research it."

	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
			openai.UserMessage(userPrompt),
		},
		Model: ChatModel,
	})
	if err != nil {
		panic(err) // Need a better error handling method here
	}
	chatResponse := chatCompletion.Choices[0].Message.Content
	UnloadModel(modelName)

	return chatResponse
}

//func callCrawlScript() (Or something along those lines)
