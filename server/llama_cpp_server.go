package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// Need to improve error handling. Right now I'm pretty sure the program just exits if it hits an error

// These probably shouldn't be hard coded. Need a .env file or a config or something like that. Probably json or toml would be best actually
// Chat model params
const ChatModelPath string = "models/Qwen2.5-7B-Instruct-Q4_K_M.gguf"
const ChatModelPort string = "8081"

// Embed model params
const EmbedModelPath string = "models/Qwen3-Embedding-8B-Q6_K.gguf"
const EmbedModelPort string = "8082"

func StartLLMServer(modelPath string, port string, embedding bool) *exec.Cmd {
	args := []string{"--model", modelPath, "--port", port}
	if embedding {
		args = append(args, "--embedding")
	}

	cmd := exec.Command("llama-server", args...)
	// Display outputs and errors in the terminal
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

// I'll double check this when I have a better understanding of Go
// Just search the internet for the relevant code...
// Need to add a timeout for this function
func WaitForServer(port string) {
	for {
		resp, err := http.Get("http://localhost:" + port + "/health") // Need to find a better way to assign the server path. Probably a .env variable or something like that
		if err == nil && resp.StatusCode == 200 {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
