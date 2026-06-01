package main

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// A bunch of variables that I'll store correctly later
const ChatModel string = "Qwen2.5-7B-Instruct-Q4_K_M"
const EmbedModel string = "Qwen3-Embedding-8B-Q5_K_M"
const ServerBaseURL string = "http://127.0.0.1:8001"
const APIKey string = "no-key"

// We'll store these elsewhere later (probably a .env file... or .env file mixed with a toml or json config file)
// Just gotta get shit working for now so we're leaving it as is

// Might be better to directly pull these values from the `curl http://localhost:8001/v1/models | jq` results instead

// Need to address the error handling everywhere. We'll leave as is for now
func main() {
	// Client creation
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	weaviateClient, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	live, err := weaviateClient.Misc().LiveChecker().Do(context.Background()) // Not really sure what this line does or where it came from. We'll figure that out later
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", live)

	// There's gotta be a way to not have to restate the client in every Weaviate related function I make dude
	// Eh you'll figure something out later, but this also might be the best way to do it
	// Client could be a global variable... but I don't like those and what happens if you need to manage multiple clients?
	// So yea... Probably leave it at that
	CreateCollection(weaviateClient, "TestCollection", "A collection to see if I can at least create an empty collection")
	var testCollection = GetCollection(weaviateClient, "TestCollection") // Don't forget to review this output to sort out the final issues with the collection config
	fmt.Printf("testCollection type: %T", testCollection)

	//result := GenerateSearchQuery(ChatModel, "What are some of the best novels of the 21st century?")
	//fmt.Println(result)
	//UnloadModel(ChatModel)
}

// Can probably just remove this one
//const WeaviateClient
