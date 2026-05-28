# What you have so far

So far I've figured out how to *properly* initiate the llama-server servers, and made a very basic single-line chat asking what the capital of france is... Not very impressive, but it's a start!

Moreover, you got the clients initiated, you found an openai Go library compatible with llama-server, you also know that there are good vector db libraries that have Go as a language for interacting with them. This begs the question, now where the hell do you start?

# Where to go next

Maybe beginning by creating a vector db the rest of the logic will come along the way. However, before you do that, you need to figure out how to incorporate that into a multi-line chat, not just a single prompt.

Actually, maybe the best way to begin is by simply figuring out how to create the vector db in Go. Then you'll actually have something to interact with.

Then, you can learn how to make a continuous chat, THEN, you can think about integrating the embedding server into said chat.

After those prerequisites, work on the crawl.py script for the web scraping fallback. After that, the rest of the logic should start falling into place.

BUT, don't forget about the need to start and stop the embedding server when you're ready. You can save a lot of RAM that way!

Sound like a plan? Yea... sounds like a plan.

# Welp... this is gonna be difficult

After getting to it, I'm running into issues with RAM (unified memory) usage when running llama-server with two models. Part of the solution can be found [here](https://www.reddit.com/r/LocalLLaMA/comments/1pmc7lk/understanding_the_new_router_mode_in_llama_cpp/) which introduces "Route Mode". Route mode essentially enables you to manage multiple AI models at the same time without restarting the server each time you switch or load a model.

There's also the issue of my limited hardware knowledge - I don't know much about how the unified memory works other than that the CPU and GPU on M chips is shared in some manner. This requires a parameter in the llama-server command to offload some of the model to the CPU (`--n-cpu-moe`). Yet, even with the offloading, the [Python version](https://github.com/crdcamp/llama-cpp-llm-embedding/blob/main/chat_demo.py) uses **significantly** less RAM than the llama-server method. This is... not great...

So, what's the solution here? Where do I go? 

I can start by studying up on unified memory and learn better how it works. That be a good place to start.

I also need to find some resources (preferably a video) for having llama-server interact with M hardware in the proper way. Otherwise, the flickering that results might cause me a seizure.

That's where we'll start then. My god am I in over my head right now.

# You're back! what did you learn?

I figured out how to use router mode! Well... by that I mean I figured out how to launch the server in router mode. I have yet to look into how the actual code can interact with it. 

Router mode is a bit new, so we'll have to make sure we're pinning the version for llama cpp somewhere somehow.

There's also the part that involves managing context. There's like a billion ways to play around with context management listed [here](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md)... and that part is gonna be wayyyy further complicated by the vector database.

That's a problem for later though. 

# Next!

Now that you got the server running in the (what I'm assuming is) the "correct" way, it's time to get this Go code initiated for interacting with the server. Refer to the openai Go docs to get a basic chat running. **Importantly**, this little demo will need to also switch between models. Just get a basic idea of that running for now.

ACTUALLY, you might want to just get the general idea of the Go code for chat completions to start, THEN make the vector database, THEN you can start switching between the models to refer to that vector database. SO, it looks like we actually need to get the crawl4ai code setup before we can begin. Either that, or just recreate the Python code using the same documents in the database. Only this time, we need a database with an SDK for Go. I think the latter would be an easier introduction.

So... **here's the plan**: 

* We need to grab those documents that were uploaded to the previous project's database and figure out how to get them into a new one that has a Go SDK.
* Then, we get to work on managing a local vector db using Go.
* The above point might require a new Go file separate from main.go. It would contain all of the logic handling for retrieving documents from, and uploading documents to the vector database.
* Then you can then start thinking about recreating the previous project in Go (as in, just have the model say "I have no idea" if the requested information isn't in the vector db)
* FINALLY, it's time to start working on crawl.py. This will be a kinda "tool" that is called in the instance the model says "I don't know"
* Then you can start doing the deep dive into managing the vector database, or as I'm referring to it: the model's "memory"

# Choosing a Vector DB

Seems like we have 2 good options: **Weaviate and Qdrant**.

Weaviate is written in Go, which is a bonus for obvious reasons. However, Qdrant is written in Rust and offers the best raw query speed.

Yet, Weaviate is better for modeling complex relationships. It's more flexible in general and "well-suited for knowledge graphs, content management systems, and applications where data relationships are as important as similarity search" [source](https://cipherprojects.com/blog/posts/weaviate-vs-qdrant-vector-database-comparison-2025/).

Qdrant takes a more minimalist approach to data modeling and has an ecosystem more focused on performance and efficiency rather than breadth of integrations.

A big plus for both is that they're open source, which for personal project purposes is a nice thing to have.

## Vector DB Conclusions

All in all, it seems like Weaviate is simpler to use and offers more capabilities out the box. I can't spend too much time learning how to use and adjust Qdrant to my purposes, as modeling a memory database (for the next project) will already be difficult enough as is.

If I had more experience and time, Qdrant would be the ideal choice due to raw performance. Yet, given that Weaviate is already written in Go and provides more tools, the choice is pretty clear.

Weavitate it is then.
