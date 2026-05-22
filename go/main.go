package main

import (
	"log"
	"os/exec"
)

var ChatModelPath string = "models/Qwen2.5-7B-Instruct-Q4_K_M.gguf"
var ChatModelPort string = "8082"

var EmbedModelPath string = "models/Qwen3-Embedding-8B-Q6_K.gguf"
var EmbedModelPort string = "8081"

func StartLLMServer(modelPath string, port string, embedding bool) {
	var startServerCmd *exec.Cmd // Declare a variable of type *exec.Cmd

	if !embedding {
		startServerCmd = exec.Command("llama-server", "--model", ModelPath, "--port", port)
	} else {
		startServerCmd = exec.Command("llama-server", "--model", ModelPath, "--port", port, "--embedding")
	}

	if err := startServerCmd.Run(); err != nil {
		log.Fatalf("Error: LLM server for model %s failed to start: %v", modelPath, err)
	}
}

func main() {
	StartLLMServer(ChatModelPath, ChatModelPort, false)
	StartLLMServer(EmbedModelPath, EmbedModelPort, true)
}
