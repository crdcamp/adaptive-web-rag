# Adaptive Web Search Rag for Local LLMs

Low parameter local LLMs are kind of dumb. So, I'm *trying* to address that problem.
# Overview

This project improves on a "rough draft" repository I wrote in Python called [llama-cpp-llm-embedding](https://github.com/crdcamp/llama-cpp-llm-embedding). The goal here is to confine a local LLM in a manner that prevents hallucinations, ensures the model is always referring to internet sources (via [crawl4ai](https://github.com/unclecode/crawl4AI)), maintains a minimal footprint on the context window, and automatically updates and checks it's "memory" using a vector database. 

All of these requirements have led to the model being **confined strictly to research-oriented questions**. The model **will not** accept any prompts other than those pertaining to research.

# Features

The features of this project are the model's ability to maintain its own research memory and its potential to run on weaker hardware (a base Macbook M4 with 16GB of RAM in my case). However, the latter point needs some work, as I need to take a lot of time after the project's completion to study [inference engineering](https://newsletter.pragmaticengineer.com/p/what-is-inference-engineering).

There are also many prompt validation checks found in the code that are conducted by the model (one of them is even a sort of decision tree... which I thought was pretty cool). This ensures that the model never strays from its intended function as a research assistant. Ask it irrelevant questions, try to manipulate it, or insult it all you'd like, the context window will reset and the model will not give in. It's for this reason that the project doesn't run as fast as I initially hoped it would on my hardware (discussed further in the [Challenges and Limitations](#challenges-and-limitations) section of this document).

# Requirements

* [**Qwen3-8B-GGUF**](https://huggingface.co/Qwen/Qwen3-8B-GGUF?show_file_info=Qwen3-8B-Q5_K_M.gguf): The LLM used in this project (quantized down to 5 bits). Save the model to the `models` directory. In a future iteration, this will automatically be handled by [HuggingFace](https://huggingface.co/). For now, you'll have to manually install the model.
* [**Docker**](https://www.docker.com/): Must be installed and running in the background. Due to my new and limited experience with Docker (and the fact my brother won't let me test the installation on his laptop), I cannot guarantee it will handle all requirements correctly.

# Usage

Using your terminal, clone the repository and `cd` into the project directory:

```bash
git clone https://github.com/crdcamp/adaptive-web-rag.git
cd adaptive-web-rag
```

To get started, run this in your terminal to setup up Docker:

```bash
docker compose build
docker compose up
```

## Example Prompts

Here are some example curl commands to send to the model when you have Docker up and running (adding an interface to the project would've just been too exhausting. I'm not much of a front-end guy):

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

(I could only think an insult for this one).

```bash
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d "{\"UserPrompt\": \"You're so useless. How are you so dumb?\"}"
```

# Project Structure

This project has three different services hosted on three different ports; all of which are constantly interacting with one another to make this thing run:

1) A [llama.cpp](https://github.com/ggml-org/llama.cpp) instance of [llama-server](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README-dev.md).
1) A [Weaviate](https://weaviate.io/) vector database.
1) A locally hosted Golang server for the chat interface. This is what accepts the above curl commands and passes them to the model.

## Prompt/Data Pipeline

**Initial prompt handling**

* User enters a prompt.
* The LLM determines the validity of the prompt in regards to whether it's a research-oriented question or not. Returns `VALID` or `INVALID`.
* If the prompt is `VALID`, continue. Otherwise, redirect the user and inform them their prompt isn't relevant to the model's capabilities.

**`VALID` prompt handling:**

* Continue processing the request by checking the vector database based on a search query generated by the model.
* Have the model validate the results from the vector database in relation to the user's prompt. Returns `RELEVANT` or `IRRELEVANT` in response to the vector database content.

**If the vector database results are `RELEVANT`:**

* If **any** of the retrieved vector database results are valid, answer the user's prompt by strictly referring to the retrieved information.
* Reset context window. Conversation history is not preserved and the LLM cannot be further prompted based on the context of the previous prompt (more on that in [Challenges and Limitations](challenges-and-limitations))

**If the vector database results are `IRRELEVANT`:**

* Conduct an internet search by reusing the generated search query and passing it to `crawl.py`.
* Determine the validity of the retrieved search results. If they are relevant, answer the user's prompt and then store the results in the vector database. Otherwise, apologize to the user and explain 

## If You Want to Experiment a Bit...

* You're welcome to experiment with different LLMs by dropping them in the `models` directory. Just be aware that they must be in a [GGUF](https://apxml.com/posts/gguf-explained-llm-file-format) format in order to work with llama.cpp and that you'll have to adjust `models.ini` and `compose.yaml` accordingly (feels good to be the guy who doesn't fully explain something in software).
* You can also adjust the system prompts for the model in the `prompts` directory if you'd like to mess around with the model's behavior a bit.

# Challenges and Limitations

# Issues to Address in the Next Iteration
