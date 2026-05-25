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

const ChatBaseURL string = "http://127.0.0.1:8001/v1"
const EmbedBaseURL string = "http://127.0.0.1:8002/v1"

func main() {
	EmbedDocuments()
}
