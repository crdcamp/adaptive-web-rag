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
