package main

import (
	"log"
	"os"
	"os/exec"
)

// Chat model params
var ChatModelPath string = "models/Qwen2.5-7B-Instruct-Q4_K_M.gguf"
var ChatModelPort string = "8082"

// Embed model params
var EmbedModelPath string = "models/Qwen3-Embedding-8B-Q6_K.gguf"
var EmbedModelPort string = "8081"

func StartLLMServer(modelPath string, port string, embedding bool) *exec.Cmd {
	args := []string{"--model", modelPath, "--port", port}
	if embedding {
		args = append(args, "--embedding")
	}

	cmd := exec.Command("llama-server", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Error: LLM server for model %s failed to start: %v", modelPath, err)
	}
	return cmd
}

func StopLLMServer(cmd *exec.Cmd) {
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		log.Fatalf("Error: Failed to stop LLM server: %v", err)
	}
	cmd.Wait()
}
