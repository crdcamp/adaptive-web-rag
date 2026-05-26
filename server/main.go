package main

// Much of the following code comes from here: https://pkg.go.dev/github.com/openai/openai-go/v3#section-readme

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// imported as openai

// Initialize a chat by initializing and connecting to the chat server
// We always want the chat server running, but we need to conditionally run and stop the embed server

// So, begin by just getting a basic chat working, then we can branch out from there
// Use the HTTP endpoints and check if there is some sort of default server chat format you should be following
// Before branching out beyond the basic chat

// After that, we can return to the embedding function and work out the logic from there

// Also... don't forget to make the embeddings and document splitting happen simultaneously

// Here are some more variables that should probably be stored differently
// (not a huge priority while I'm learning an entire new language though)
const ChatBaseURL string = "http://127.0.0.1:8001/v1"
const EmbedBaseURL string = "http://127.0.0.1:8002/v1"

func main() {
	ctx := context.Background()
	chatClient := openai.NewClient(
		option.WithBaseURL(ChatBaseURL),
		option.WithAPIKey("no-key"),
	)
	embedClient := openai.NewClient(
		option.WithBaseURL(EmbedBaseURL),
		option.WithAPIKey("no-key"),
	)
}

// Here's how you do it with the openai client:

// You might wanna "pin" the version with this:
// go get -u 'github.com/openai/openai-go/v3@v3.37.0'

// package main

// import (
//     "github.com/openai/openai-go/v3"
//     "github.com/openai/openai-go/v3/option"
// )

// var (
//     chatClient      = openai.NewClient(
//         option.WithBaseURL("http://localhost:8001/v1"),
//         option.WithAPIKey("no-key"),
//     )
//     embeddingClient = openai.NewClient(
//         option.WithBaseURL("http://localhost:8002/v1"),
//         option.WithAPIKey("no-key"),
//     )
// )

// Then to use them:

// // Chat
// chatClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{...})

// // Embeddings
// embeddingClient.Embeddings.New(ctx, openai.EmbeddingNewParams{...})
