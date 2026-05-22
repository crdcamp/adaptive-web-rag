#!/usr/bin/env bash

# Load models
llama-server --model models/Qwen3-Embedding-8B-Q6_K.gguf --port 8081 --embedding
llama-server --model models/Qwen2.5-7B-Instruct-Q4_K_M.gguf --port 8082
