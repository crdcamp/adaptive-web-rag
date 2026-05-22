package main

import (
	"os/exec"
)

// Initialize chat and embedding model servers
func main() {
	chatServerCmd, chatErr := exec.Command("llama-server", "--model", "models/Qwen2.5-7B-Instruct-Q4_K_M.gguf", "--port", "8082")
	embedServerCmd, embedErr := exec.Command("llama-server", "--model", "models/Qwen3-Embedding-8B-Q6_K.gguf", "--port", "8081", "--embedding")
	if chatErr || embedErr != nil {
		panic(err)
	}
}
