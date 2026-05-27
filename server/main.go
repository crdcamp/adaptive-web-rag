package main

// We'll store these elsewhere later (probably a .env file... or .env file mixed with a toml or json config file)
// Just gotta get shit working for now so we're leaving it as is
const ChatModel string = "Qwen2.5-7B-Instruct-Q4_K_M"
const ServerBaseURL string = "http://127.0.0.1:8001"
const APIKey string = "no-key"

// Main function is completely useless other than for testing as of now
func main() {
	GenerateSearchQueries("Tell me about vector databases")
}
