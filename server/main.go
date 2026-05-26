package main

// Much of the following code comes from here: https://pkg.go.dev/github.com/openai/openai-go/v3#section-readme

import (
	"context"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
)

const ChatBaseURL string = "http://127.0.0.1:8001/v1"
const EmbedBaseURL string = "http://127.0.0.1:8002/v1"

// For later
// embedClient := openai.NewClient(
	// 	option.WithBaseURL(EmbedBaseURL),
	// 	option.WithAPIKey("no-key"),
	// )

func main() {
	ctx := context.Background()
	chatClient := openai.NewClient(
		option.WithBaseURL(ChatBaseURL),
		option.WithAPIKey("no-key"),
	)

	question := "What is the capital of France?"

	resp, err := chatClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "models/Qwen2.5-7B-Instruct-Q4_K_M.gguf",
		Messages: []openai.ChatCompletionMessageParamUnion{
            openai.UserMessage(question),
        },
    })
    if err != nil {
        panic(err)
    }
    println(resp.Choices[0].Message.Content)
}
