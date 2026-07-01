package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go/v3"
)

// Unload a model from memory using the `/models/unload` HTTP endpoint.
// Available models and their status can be displayed using `curl http://localhost:8080/v1/models | jq`
func UnloadModel(modelName string) {
	var unloadURL = AppConfig.LlamaServer + "/models/unload"
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
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Unload failed: %s - %s\n", resp.Status, string(body))
		return
	}
	defer resp.Body.Close()
	//fmt.Printf("Status: %s\n", resp.Status)
	fmt.Println("Model unloaded:", modelName)
}

func CreateChatCompletion(client openai.Client, modelName, userPrompt string, systemPrompt string) {
	ctx := context.Background()

	fmt.Printf("Creating chat completion\nUser prompt: %q\nSystem prompt: %q\n", userPrompt, systemPrompt)
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model: modelName,
	})
	if err != nil {
		panic(err)
	}
	chatResponse := chatCompletion.Choices[0].Message.Content
	fmt.Printf("Chat completion response: %q\n", chatResponse)
}

// Generate one internet search query and save the result to `server/crawl_data/user_prompt.md`.
// The resulting query may be used by `crawl.py`.
func GenerateSearchQuery(client openai.Client, modelName string, userPrompt string) {
	ctx := context.Background()

	fmt.Printf("Generating search query for user prompt: %q\n", userPrompt)
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

	chatResponseByte := []byte(strings.Trim(chatResponse, `"`))
	path := filepath.Join("crawl_data/user_prompt.md")

	// This needs a check for if the file exists: https://golangtutorial.dev/tips/check-if-a-file-exists-or-not-in-go/
	err = os.WriteFile(path, chatResponseByte, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Search query saved to `server/crawl_data/user_prompt.md`")
}

// Calls crawl.py to conduct web search. Results are saved to `server/crawl_data/crawl_results.json`.
func CallCrawlScript() {
	cmd := exec.Command("python3", "crawl.py")

	// Output to terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

//func RefineVectorDBSearch() {}
