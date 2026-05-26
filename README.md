# Adaptive Web Search Rag for Local LLMs

In this project, I learn a new coding language (Golang), start hosting My LLMs locally with HTTP endpoints (rather than running them within a Python script), and build off of the [llama-cpp-llm-embedding](https://github.com/crdcamp/llama-cpp-llm-embedding) repository by adding a new web search fallback capability to the embedding project.

This is ultimately meant to run on a Macbook base M4 chip, but running these models using `llama-server` instead of `llama-cpp-python` is showing to be a bit more involved than I initially thought.

I'm essentially rewriting the entire previous repository in a new language I know nothing about (other than that it scales very well) as well, so there are gonna be some challenges along the way. This one might take a while to get off the ground.

Keep in mind, this project is in its **very** early stages. Don't look at me!!!

Here's the installation setup so far:

# Installation

```bash
# Download source code
git clone https://github.com/crdcamp/adaptive-web-rag.git
cd adaptive-web-rag

# Install llama.cpp
brew install llama.cpp

# Install openai Go package
cd server
go get -u 'github.com/openai/openai-go/v3@v3.37.0'
cd ..

# Setup Python environment and crawl4ai
python3 -m venv
source venv/bin/activate
pip install -r requirements.txt
crawl4ai-setup
```

## Installing Models

For simplicity, I've been manually installing models. Eventually I'll start interacting with the HuggingFace CLI to streamline the installation a bit. For now, just visit these web pages and download the models into the `models` directory:

- [Qwen/Qwen3-Embedding-8B-GGUF-Q5_K_M](https://huggingface.co/Qwen/Qwen3-Embedding-8B-GGUF?show_file_info=Qwen3-Embedding-8B-Q5_K_M.gguf)
- [bartowski/Qwen2.5-7B-Instruct-Q4_K_M.gguf](https://huggingface.co/bartowski/Qwen2.5-7B-Instruct-GGUF/blob/main/Qwen2.5-7B-Instruct-Q4_K_M.gguf)

# Running the Servers

If you're on a Mac base M4 chip like me... DO NOT run these commands simultaneously! If you do, get ready for a lot of screen flickering and the need to force shutdown your computer.

To run the **chat** server:

```bash
llama-server -m models/Qwen2.5-7B-Instruct-Q4_K_M.gguf --port 8001 -c 1024
```

To run the **embed** server:

```bash
llama-server -m models/Qwen3-Embedding-8B-Q6_K.gguf --port 8002 -c 1024
```

# Fixing llama-server RAM issues with new commands

This [GitHub source](https://github.com/ggml-org/llama.cpp/discussions/15396) serves as a good introduction for managing memory usage on Macbooks in general. However, I can't seem to get this RAM (or "unified memory") usage down to what the Python script was able to!

Start **instruct** sever:

```bash
llama-server -m models/Qwen2.5-7B-Instruct-Q4_K_M.gguf --n-cpu-moe 12 -c 2048 --port 8001
```

Start **embedding** server:

```bash
llama-server -m models/Qwen3-Embedding-8B-Q6_K.gguf --n-cpu-moe 12 -c 2048 --port 8002
```

* `--n-cpu-moe`: Number of MoE layers N to keep on the CPU. This is used in hardware configs that cannot fit the models fully on the GPU. The specific value depends on your memory resources and finding the optimal value requires some experimentation

* `-c`: Specify the context size to use. More context requires more memory. Both gpt-oss models have a maximum context of 128k tokens. Use -c 0 to set to the model's default

* `--no-mmap`: Disables memory-mapping when loading the model file
