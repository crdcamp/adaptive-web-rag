package main

// OpenAI API reference: https://pkg.go.dev/github.com/openai/openai-go/v3#section-readme

// INSTEAD OF SAVING THE JSON FILE,
// POST IT WITH HTTP AND THEN ACCESS THAT FROM THE
// PYTHON FILE

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
)

// This entire function is AI slop... will review when I have more understanding of Go
// I have a feeling (I wonder why?) that this function could be further simplified
func UnloadModel(modelID string) error {
	body, err := json.Marshal(map[string]string{"model": modelID})
	if err != nil {
		return fmt.Errorf("failed to marshal unload request: %w", err)
	}

	resp, err := http.Post(
		ServerBaseURL+"/models/unload",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("unload request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unload returned status %d", resp.StatusCode)
	}
	return nil
}

// Given the user's prompt, generate 5 search queries to be used by `crawl.py` for web scraping
func GenerateSearchQueries(userPrompt string) {
	ctx := context.Background()
	chatClient := openai.NewClient(
		option.WithBaseURL(ServerBaseURL),
		option.WithAPIKey(APIKey),
	)

	systemPrompt := `You are a search query generator. When given a question or topic, generate exactly five concise search engine queries a person could enter into a browser to research it.

Respond ONLY with a JSON object in this exact format, no other text:
{"queries": ["query1", "query2", "query3", "query4", "query5"]}`

	resp, err := chatClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: ChatModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
	})
	if err != nil {
		panic(err) // Need to double check if `panic` is a good idea here
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
	println(string(out))
	UnloadModel(ChatModel)
}

//func callCrawlScript()
