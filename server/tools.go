package main

import (
	"context"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
)

func GenerateSearchQueries(userPrompt string) {
	ctx := context.Background()
	chatClient := openai.NewClient(
		option.WithBaseURL(BaseURL),
		option.WithAPIKey(APIKey),
	)

	systemPrompt := "You are a search query generator. When given a question or topic, generate exactly five concise search engine queries a person could enter into a browser to research it."

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
	println(resp.Choices[0].Message.Content)
}

//func callCrawlScript()
