# Adaptive Web Search Rag for Local LLMs

In this project, I learn a new coding language (Golang), start hosting My LLMs locally with HTTP endpoints (rather than running them within a Python script), learn Docker (huge pain in the butt), and build off of the [llama-cpp-llm-embedding](https://github.com/crdcamp/llama-cpp-llm-embedding) repository by adding a new web search fallback capability to the embedding project. Ignorantly configuring Docker and the URLs was also a huge pain point for this one (for what it's worth).

Anyway, welcome to my favorite project I've done to date where I create a "proof of concept" for a locally hosted AI involving a vector database based on web search results that the program generates and embeds on it's own.


I'm essentially rewriting the entire [previous repository](https://github.com/crdcamp/llama-cpp-llm-embedding) in a new language (Golang) I know nothing about (other than that it scales very well). So, this one's a bit rough around the edges (I'm mostly referring Golang to error handling here).

This project is designed to run on a Macbook base M4 chip at least... decently.

## The Model I'm Using

For simplicity, I've been manually installing models. Eventually I'll start interacting with the HuggingFace CLI to streamline the installation a bit. For now, just visit these web pages and download the models into the `models` directory:

- [Qwen3-8B-GGUF](https://huggingface.co/Qwen/Qwen3-8B-GGUF?show_file_info=Qwen3-8B-Q5_K_M.gguf)

## Issues to address in future projects

* Embedding and Instruct model are not compatible (got in too far to address this).
* Doesn't check whether site had been visited when embedding.
* Needs to check if question requires a search result in the first place.
* Needs a way to adjust desired search results from Go scripts.
* Web search data quality isn't exactly what we need (could be solved with chain-of-thought web search)
* Web search doesn't produce contextual searches well at all. Could also be solved with the above suggestion.
* `docker-compose.yml` doesn't set everything up in a fresh environment.
* Data handling between Go and Python could be much more elegant (maybe they can share data with memory?).
* No updates on progress sent from server.

# To Do
- [ ] Add python environment creation and `pip install -r requirements.txt` to `Dockerfile`.
- [ ] Add `\n` to create chat completion output (or write a function for the write bytes that includes `\n`).

# Starting Docker

Easy:

```bash
docker compose up
```

To show a list of the models discovered by the Router:

```bash
curl http://localhost:8080/v1/models | jq
```

```bash
docker compose build
```

# Some Test Curl Commands for Ya

## Declarative Questions (definitely not the correct terminology)

```bash
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d '{"UserPrompt": "Find out what the most popular restaraunts in Genieva Switzerland are"}'
```

## Regular Questions

```bash
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d '{"UserPrompt": "What are some of the greatest mysteries throughout ancient history?"}'
```

```bash
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d '{"UserPrompt": "What are some of the greatest mysteries throughout ancient history?"}'
```

```bash
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d '{"UserPrompt": "How old would the founder of the company that acquired Instagram have been when Instagram was founded?"}'
```
