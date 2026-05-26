# Adaptive Web Search Rag for Local LLMs

In this project, I learn a new coding language (Golang), start hosting My LLMs locally (rather than running them within a Python script), and build off of the [llama-cpp-llm-embedding](https://github.com/crdcamp/llama-cpp-llm-embedding) repository by adding a new web search fallback capability to the embedding project.

I'm essentially rewriting the entire previous repository in a new language I know nothing about (other than that it scales very well), so there are gonna be some challenges along the way. This one might take a while to get off the ground.

Keep in mind, this project is in its **very** early stages. Don't look at me!

Here's the installation setup so far:

```bash
# Install llama.cpp
brew install llama.cpp

# Setup Python environment and crawl4ai
python3 -m venv
source venv/bin/activate
pip install -r requirements.txt
crawl4ai-setup
```

# Run the Servers

To run the **chat** server:

```bash
llama-server -m models/Qwen3-Embedding-8B-Q6_K.gguf --port 8001
```

To run the **embed** server:

```bash
llama-server -m models/Qwen3-Embedding-8B-Q6_K.gguf --port 8002
```
