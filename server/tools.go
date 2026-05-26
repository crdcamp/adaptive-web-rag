package main

import (
	"context"
	"encoding/json"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
)

func GenerateSearchQueries(userPrompt string) {
	ctx := context.Background()
	chatClient := openai.NewClient(
		option.WithBaseURL(BaseURL),
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
		panic(err)
	}

	raw := resp.Choices[0].Message.Content

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
}

//func callCrawlScript()
