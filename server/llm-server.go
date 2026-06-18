package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go/v3"
)

// Load a model into memory using the `/models/load` HTTP endpoint.
// Available models and their status can be displayed using `curl http://localhost:8080/v1/models | jq`
func LoadModel(modelName string) {
	var loadURL = LlamaBaseUrl + "/models/load"
	payload, err := json.Marshal(map[string]string{"model": modelName})
	if err != nil {
		panic(err)
	}
	fmt.Println("Loading model:", modelName)
	resp, err := http.Post(loadURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Model unloaded:", modelName)
}

// Unload a model from memory using the `/models/unload` HTTP endpoint.
// Available models and their status can be displayed using `curl http://localhost:8080/v1/models | jq`
func UnloadModel(modelName string) {
	var unloadURL = LlamaBaseUrl + "/v1/models/unload"
	// Need to research more into json encoding in Go. I have no idea how this works at the moment
	payload, err := json.Marshal(map[string]string{"model": modelName})
	if err != nil {
		panic(err)
	}
	fmt.Println("Unloading model:", modelName)
	resp, err := http.Post(unloadURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err) // Need a better error handling method here
	}
	defer resp.Body.Close()
	//fmt.Printf("Status: %s\n", resp.Status)
	fmt.Println("Model unloaded:", modelName)
}

// Generate one internet search query and save the result to `server/crawl_data/user_prompt.md`.
// The resulting query may be used by `crawl.py`.
func GenerateSearchQuery(client openai.Client, modelName string, userPrompt string) {
	ctx := context.Background()

	fmt.Println("Generating search query for user prompt:", userPrompt)
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a search query generator. When given a question or topic, generate ONE search engine query that a person could enter into a browser to research it."),
			openai.UserMessage(userPrompt),
		},
		Model: modelName,
	})
	if err != nil {
		panic(err)
	}
	chatResponse := chatCompletion.Choices[0].Message.Content
	fmt.Println("Search query generated:", chatResponse)

	fmt.Println("Saving prompt to `server/crawl_data/user_prompt.md`")
	chatResponseByte := []byte(strings.Trim(chatResponse, `"`))
	path := filepath.Join("crawl_data/user_prompt.md")
	err = os.WriteFile(path, chatResponseByte, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Prompt saved to `server/crawl_data/user_prompt.md`")
}
