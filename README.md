# Adaptive Web Search Rag for Local LLMs

In this project, I learn a new coding language (Golang), start hosting My LLMs locally with HTTP endpoints (rather than running them within a Python script), learn Docker (huge pain in the butt), and build off of the [llama-cpp-llm-embedding](https://github.com/crdcamp/llama-cpp-llm-embedding) repository by adding a new web search fallback capability to the embedding project. Ignorantly configuring Docker and the URLs was also a huge pain point for this one (for what it's worth).

Anyway, welcome to my favorite project I've done to date where I create a "proof of concept" for a locally hosted AI involving a vector database based on web search results that the program generates and embeds on it's own.


I'm essentially rewriting the entire [previous repository](https://github.com/crdcamp/llama-cpp-llm-embedding) in a new language (Golang) I know nothing about (other than that it scales very well). So, this one's a bit rough around the edges (I'm mostly referring Golang to error handling here).

This project is designed to run on a Macbook base M4 chip, since I'm a firm believer that when AI companies need to start making profit the current token cost structure is going to be f****d and unaffordable. Oh how I despise technical debt even with 0 professional experience!

## The Models I'm Using

For simplicity, I've been manually installing models. Eventually I'll start interacting with the HuggingFace CLI to streamline the installation a bit. For now, just visit these web pages and download the models into the `models` directory:

- [Qwen/Qwen3-Embedding-8B-GGUF-Q5_K_M](https://huggingface.co/Qwen/Qwen3-Embedding-8B-GGUF?show_file_info=Qwen3-Embedding-8B-Q5_K_M.gguf)
- [Qwen/Qwen2.5-7B-Instruct-GGUF](https://huggingface.co/Qwen/Qwen2.5-7B-Instruct-GGUF?show_file_info=qwen2.5-7b-instruct-q4_k_m-00001-of-00002.gguf)

## Issues to address in future projects

* Doesn't check whether site had been visited when embedding.
* Needs to check if question requires a search result in the first place.
* Needs a way to adjust desired search results from Go scripts.
* Web search data quality isn't exactly what we need (could be solved with chain-of-thought web search)
* Web search doesn't produce contextual searches well at all. Could also be solved with the above suggestion.
* `docker-compose.yml` doesn't set everything up in a fresh environment.
* Data handling between Go and Python could be much more elegant (maybe they can share data with memory?).

# Starting Docker

Easy:

```bash
docker compose up
```

To show a list of the models discovered by the Router:

```bash
curl http://localhost:8080/v1/models | jq
```
