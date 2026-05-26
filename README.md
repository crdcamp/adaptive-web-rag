# Adaptive Web Search Rag for Local LLMs

In this project, I learn a new coding language (Golang), start hosting My LLMs locally with HTTP endpoints (rather than running them within a Python script), and build off of the [llama-cpp-llm-embedding](https://github.com/crdcamp/llama-cpp-llm-embedding) repository by adding a new web search fallback capability to the embedding project.

I'm essentially rewriting the entire previous repository in a new language I know nothing about (other than that it scales very well), so there are gonna be some challenges along the way. This one might take a while to get off the ground.

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

- [Qwen/Qwen3-Embedding-8B-GGUF](https://huggingface.co/Qwen/Qwen3-Embedding-8B-GGUF?show_file_info=Qwen3-Embedding-8B-Q6_K.gguf)
- [bartowski/Qwen2.5-7B-Instruct-Q4_K_M.gguf](https://huggingface.co/bartowski/Qwen2.5-7B-Instruct-GGUF/blob/main/Qwen2.5-7B-Instruct-Q4_K_M.gguf)

# Running the Servers

To run the **chat** server:

```bash
llama-server -m models/Qwen3-Embedding-8B-Q6_K.gguf --port 8001
```

To run the **embed** server:

```bash
llama-server -m models/Qwen3-Embedding-8B-Q6_K.gguf --port 8002
```
