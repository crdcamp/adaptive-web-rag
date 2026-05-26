package main

// Much of the following code comes from here: https://pkg.go.dev/github.com/openai/openai-go/v3#section-readme

// We'll store these elsewhere later
const ChatModel string = "Qwen2.5-7B-Instruct-Q4_K_M"
const ChatBaseURL string = "http://127.0.0.1:8001/v1"
const EmbedBaseURL string = "http://127.0.0.1:8002/v1"
const APIKey string = "no-key"


// For later
// embedClient := openai.NewClient(
// 	option.WithBaseURL(EmbedBaseURL),
// 	option.WithAPIKey("no-key"),
// )

func main() {
	GenerateSearchQueries("Tell me about vector databases")
}
