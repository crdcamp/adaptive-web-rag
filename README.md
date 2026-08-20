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
* Determine the validity of the retrieved search results. If they are relevant, answer the user's prompt and then store the results in the vector database. Otherwise, apologize to the user and explain that this project has a lot of limitations.

## If You Want to Experiment a Bit...

* You're welcome to experiment with different LLMs by dropping them in the `models` directory. Just be aware that they must be in a [GGUF](https://apxml.com/posts/gguf-explained-llm-file-format) format in order to work with llama.cpp and that you'll have to adjust `models.ini` and `compose.yaml` accordingly (feels good to be the guy who doesn't fully explain something in software).
* You can also adjust the system prompts for the model in the `prompts` directory if you'd like to mess around with the model's behavior a bit.

# Challenges and Limitations

## Challenges

This was easily the most difficult and complex project I've ever done so far. I grew tired of the abstractions of Python, wanted a more efficient language, and wanted a much greater challenge than my previous projects have posed.

This led me to Golang; a language I haven't used before and knew very little about. Between Golang's efficiency, scalability, general design towards network handling, and its small step towards a "lower-level" language, it turned out to be the perfect choice. I truly love this language and hope I get to use it professionally one day.

There was also the issue of diving head first into Docker for the first time; a software that still seems like a black box to me, but I'm glad I learned a bit about it. It was a **massive** headache, and that's exactly what I was hoping for. I'll likely continue to use and learn more about this software as I keep adding repositories. It's undoubtedly something I'll need to be comfortable with professionally one day anyway.

Moreover, due to the complexity, I had to refactor the code at least five separate times. While at times a bit exhausting, this was also sometimes pretty fun to do. Learning about what it takes to start with a blank canvas like this, iterate over and over, and see an idea come to fruition like this is a rewarding process at the end of the day. I have absolute control over every aspect of these files... and something about that makes me a happy boy. While I could take some more time to clean it up, eventually you just have to tell yourself that you need to get it done.

All in all, this has been quite the process. It's definitely increased my programming confidence and might've even eliminated that little bit of imposter syndrome I had left before beginning it.

## Limitations

Due to my goal of actually finishing this project, it comes with a lot of limitations.

### Chat Limitations

* No chat history or continuation of conversation. This would be another project in itself.
* No timestamps for chat history.
* No handling for follow up questions.
* No updates from the chat server to the user. User must watch the Docker logs to find out what's going on.

### Web Search Limitations

* No check for already visited websites.
* DockDuckGo search in Python is inconsistent in its results. The same search query can result in different URLs. Need to migrate to a different method.
* No list of blocked sites that may contain unreliable information.
* A sort of "chain of thought" operation (or something similar) could address these issues, but that's a project in itself.

### Vector Database Limitations

* No dates included with web search results. Prevents the ability to update the vector database if the information is considered outdated.
* No "chat memory" collection included. Was originally a project requirement but would take too much additional time to implement.
* No endpoint to easily manage the vector database's collections. Would require yet another Golang file with a bunch of hand-written functions to implement.

### Structural

* Good structure in general (sorta), yet the Golang scripts would ideally be but in their own directory. They would also ideally be further refactored by their functionality, yet this requires reconfiguring Docker a bit (I already refactored enough times).
* Error handling is inconsistent and not very harmonious with how the chat server operates (a byproduct of starting with a blank canvas).
* Some LLM functions send posts to the server, some don't. This also ties into the chat limitation on not sending updates. Would unfortunately require a lot of effort to implement (there were already enough concepts I had to learn in the first place).
* Data sharing from Golang scripts to `chat.py` could be handled more elegantly. As of now data is shared by reading/writing files.
