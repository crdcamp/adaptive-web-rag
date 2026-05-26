#!/usr/bin/env bash

# Start chat server
llama-server -m models/Qwen2.5-7B-Instruct-Q4_K_M.gguf --port 8001

# Start embed server
llama-server -m models/Qwen3-Embedding-8B-Q6_K.gguf --port 8002
