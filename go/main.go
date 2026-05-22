package main

import (
	"log"
	"os/exec"
)

var chatModelPath string = "models/Qwen2.5-7B-Instruct-Q4_K_M.gguf"
var embedModelPath string = "models/models/Qwen3-Embedding-8B-Q6_K.gguf"

func startServers(chatModelPath string, chatPort string, embedModelPath string, embedPort string) {
	// Initialize chat and embedding model servers
	chatServerCmd := exec.Command("llama-server", "--model", chatModelPath, "--port", chatPort)
	embedServerCmd := exec.Command("llama-server", "--model", embedModelPath, "--port", "8081", embedPort)

	if err := chatServerCmd.Run(); err != nil {
		log.Fatalf("chat server failed: %v", err)
	}
	if err := embedServerCmd.Run(); err != nil {
		log.Fatalf("embed server failed: %v", err)
	}

}

func main() {
	startServers("8082", "8081")
}
