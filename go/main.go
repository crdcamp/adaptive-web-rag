package main

import (
	"log"
	"os/exec"
	"time"
)

// Chat model params
var ChatModelPath string = "models/Qwen2.5-7B-Instruct-Q4_K_M.gguf"
var ChatModelPort string = "8082"

// Embed model params
var EmbedModelPath string = "models/Qwen3-Embedding-8B-Q6_K.gguf"
var EmbedModelPort string = "8081"

func StartLLMServer(modelPath string, port string, embedding bool) *exec.Cmd {
	var startServerCmd *exec.Cmd // Declare a variable of type *exec.Cmd

	if !embedding {
		startServerCmd = exec.Command("llama-server", "--model", modelPath, "--port", port)
	} else {
		startServerCmd = exec.Command("llama-server", "--model", modelPath, "--port", port, "--embedding")
	}

	if err := startServerCmd.Start(); err != nil {
		log.Fatalf("Error: LLM server for model %s failed to start: %v", modelPath, err)
	}
	return startServerCmd
}

func main() {
	chatCmd := StartLLMServer(ChatModelPath, ChatModelPort, false)
	embedCmd := StartLLMServer(EmbedModelPath, EmbedModelPort, true)

	// For testing: kill the server after one minute
	time.Sleep(1 * time.Minute)

	chatCmd.Process.Kill()
	embedCmd.Process.Kill()
}
