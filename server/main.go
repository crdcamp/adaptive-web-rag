package main

// Initialize a chat by initializing and connecting to the chat server
// We always want the chat server running, but we need to conditionally run and stop the embed server

// So, begin by just getting a basic chat working, then we can branch out from there
// Use the HTTP endpoints and check if there is some sort of default server chat format you should be following
// Before branching out beyond the basic chat

// After that, we can return to the embedding function and work out the logic from there

// Also... don't forget to make the embeddings and document splitting happen simultaneously

// Here are some more variables that should probably be stored differently
// (not a huge priority while I'm learning an entire new language though)

func main() {
	EmbedDocuments()
}

// Here's how you do it with the openai client:

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
