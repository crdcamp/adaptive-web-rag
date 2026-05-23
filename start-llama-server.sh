#!/usr/bin/env bash

# Start servers in background
llama-server --model models/Qwen2.5-7B-Instruct-Q4_K_M.gguf --port 8082 &
llama-server --model models/Qwen3-Embedding-8B-Q6_K.gguf --port 8081 --embedding &

# Wait for them to be ready (they take a few seconds to load)
sleep 10

# Run the Go app
cd go && go run .
