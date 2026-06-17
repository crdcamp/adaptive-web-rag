package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/v3"
)

func UnloadModelRw(modelName string) {
	const unloadURL = LlamaBaseUrl + "/v1/models/unload"
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

func GenerateSearchQueryRw(client openai.Client, modelName string, userPrompt string) {
	ctx := context.Background()

	fmt.Printf("Generating search query for user prompt: %q\n", userPrompt)
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a search query generator. When given a question or topic, generate a search engine query that a person could enter into a browser to research it."),
			openai.UserMessage(userPrompt),
		},
		Model: modelName,
	})
	if err != nil {
		panic(err)
	}
	chatResponse := chatCompletion.Choices[0].Message.Content
	fmt.Printf("Search query generated: %q", chatResponse)
	UnloadModelRw(ChatModel)
}
