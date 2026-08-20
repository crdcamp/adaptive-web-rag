# Adaptive Web Search Rag for Local LLMs

Low parameter local LLMs are kind of dumb. So, I'm *trying* to address that problem.

# Overview

This project improves on a "rough draft" repository I wrote in Python called [llama-cpp-llm-embedding](https://github.com/crdcamp/llama-cpp-llm-embedding). The goal here is to confine a local LLM in a manner that prevents hallucinations, ensures the model is always referring to internet sources (via [crawl4ai](https://github.com/unclecode/crawl4AI)), maintains a minimal footprint on the context window, and automatically updates and checks it's "memory" using a vector database. 

All of these requirements have led to the model being confined strictly to research-oriented questions. The model **will not** accept any prompts other than those pertaining to research.

# Features

The features of this project are the model's ability to maintain its own research memory and its potential to run on weaker hardware (a base Macbook M4 with 16GB of RAM in my case). However, the latter point needs some work, as I need to take a lot of time after the project's completion to study [inference engineering](https://newsletter.pragmaticengineer.com/p/what-is-inference-engineering).

There are also many prompt validation checks found in the code that are conducted by the model (one of them is even a sort of decision tree... which I thought was pretty cool). This ensures that the model never strays away from it's intended function as a research assistant. Ask it irrelevant questions, try to manipulate it, or insult it all you'd like, the context window will reset and the model will not give in. It's for this reason that the project doesn't run as fast as I initially hoped it would on my hardware (discussed further in the [Challenges and Limitations](#challenges-and-limitations) section of this document).

# Requirements

* [Docker](https://www.docker.com/): Must be installed and running in the background. Due to my new and limited experience with Docker (and the fact my brother won't let me test the installation on his laptop), I cannot guarantee it will handle all requirements correctly.

# Usage

Using your terminal, clone the repository and `cd` into the directory:

```bash
git clone https://github.com/crdcamp/adaptive-web-rag.git
cd adaptive-web-rag
```

To get started, run this in your terminal as well:

```bash
docker compose build
docker compose up
```

## Prompt Examples

Here are some example curl requests to send to the model when you have Docker up and running 

(Adding an interface to the project would've just been too exhausting. I'm not much of a front-end guy).

### Intended Questions

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

### Declarative Questions

```bash
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d '{"UserPrompt": "Find out what the most popular restaraunts in Genieva Switzerland are"}'
```

```bash
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d '{"UserPrompt": "Tell me about the benefits and tradeoffs of a vector database"}'
```

### Out-Of-Scope Questions

(I could only think of insults for this one).

```bash
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d "{\"UserPrompt\": \"You're so useless. How are you so dumb?\"}"
```

# Project Structure

# The Model I'm Using

# Challenges and Limitations

# Issues to Address in the Next Iteration
