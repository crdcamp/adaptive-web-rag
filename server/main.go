package main

import "fmt"

// We'll store these elsewhere later (probably a .env file... or .env file mixed with a toml or json config file)
// Just gotta get shit working for now so we're leaving it as is

// Might be better to directly pull these values from the `curl http://localhost:8001/v1/models | jq` results instead

const ChatModel string = "Qwen2.5-7B-Instruct-Q4_K_M"
const EmbedModel string = "Qwen3-Embedding-8B-Q5_K_M"
const ServerBaseURL string = "http://127.0.0.1:8001"
const APIKey string = "no-key" // Can probably just remove this one

// Main function is completely useless other than for testing as of now
func main() {
	result := GenerateSearchQueries("What are some of the best novels of the 21st century?")
	fmt.Println(result)
	UnloadModel(ChatModel)
}
