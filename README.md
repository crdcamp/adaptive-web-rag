# Adaptive Web Search Rag for Local LLMs

**This entire doc is irrelevant at the moment**

In this project, I learn a new coding language (Golang), start hosting My LLMs locally with HTTP endpoints (rather than running them within a Python script), and build off of the [llama-cpp-llm-embedding](https://github.com/crdcamp/llama-cpp-llm-embedding) repository by adding a new web search fallback capability to the embedding project.

This is ultimately meant to run on a Macbook base M4 chip, but running these models using `llama-server` instead of `llama-cpp-python` is showing to be a bit more involved than I initially thought.

I'm essentially rewriting the entire previous repository in a new language I know nothing about (other than that it scales very well) as well, so there are gonna be some challenges along the way. This one might take a while to get off the ground.

Keep in mind, this project is in its **very** early stages. Don't look at me!!!

Here's the installation setup so far:

# Installation (this is far from accurate)

```bash
# Install llama.cpp
brew install llama.cpp

# Download source code
git clone https://github.com/crdcamp/adaptive-web-rag.git
cd adaptive-web-rag
mkdir models

# Setup Python environment and crawl4ai
python3 -m venv
source venv/bin/activate
pip install -r requirements.txt
crawl4ai-setup

# Install openai Go package
cd server
go get -u 'github.com/openai/openai-go/v3@v3.37.0'
go mod init github.com/weaviate-go-client
go mod tidy

# Return to the project root
cd ..
```

## Installing Models

For simplicity, I've been manually installing models. Eventually I'll start interacting with the HuggingFace CLI to streamline the installation a bit. For now, just visit these web pages and download the models into the `models` directory:

- [Qwen/Qwen3-Embedding-8B-GGUF-Q5_K_M](https://huggingface.co/Qwen/Qwen3-Embedding-8B-GGUF?show_file_info=Qwen3-Embedding-8B-Q5_K_M.gguf)
- [Qwen/Qwen2.5-7B-Instruct-GGUF](https://huggingface.co/Qwen/Qwen2.5-7B-Instruct-GGUF?show_file_info=qwen2.5-7b-instruct-q4_k_m-00001-of-00002.gguf)

# Starting Docker

Easy:

```bash
docker compose up
```

# Starting the Server

* [Server troubleshooting guide](https://www.youtube.com/watch?v=1_L9cG-X2eY)
* [Managing hardware interaction for M chips](https://github.com/ggml-org/llama.cpp/discussions/15396)

The following command starts the server in router mode. The server is set to only allow one model to be loaded at a time. More information on the parameters can be found [here](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md).

* You can visit the local host in your browser and be greeted with a web UI. However, if you want to not use the web UI, add this parameter: `--no-webui`.
* We're gonna use the `--verbose` parameter for now as well.
* You also might wanna change the `--models-autoload` parameter to `--no-models-autoload` later.
* There's also the issue of whether or not to use a jinja template. Investigate that parameter later.

```bash
llama-server --models-dir models/ --n-cpu-moe 12 --mlock -c 2048 --verbose --models-max 1 --models-autoload --port 8001
```

* `--n-cpu-moe`: Number of MoE layers N to keep on the CPU. This is used in hardware configs that cannot fit the models fully on the GPU. The specific value depends on your memory resources and finding the optimal value requires some experimentation
* `-c`: Specify the context size to use. More context requires more memory. Both gpt-oss models have a maximum context of 128k tokens. Use -c 0 to set to the model's default
* `--no-mmap`: Disables memory-mapping when loading the model file
* * `--mlock`: Forces the system to keep model in RAM rather than swapping or compressing.


To show a list of the models discovered by the Router:

```bash
curl http://localhost:8001/v1/models | jq
```
