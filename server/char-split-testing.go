package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/crdcamp/charsplitter"
)

func TestSplit() {
	// Read `crawl_results.json`
	//ctx := context.Background()
	content, err := os.ReadFile("crawl_data/crawl_results.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Create a struct for each key's values
	type HrefContent struct {
		Href    string `json:"href"`
		Content string `json:"content"`
	}

	// Unmarshal the data into HrefContent
	jsonMap := map[string]HrefContent{}
	err = json.Unmarshal(content, &jsonMap)
	if err != nil {
		panic(err)
	}

	// Text splitter
	splitter := charsplitter.New(
		charsplitter.WithChunkSize(512),
		charsplitter.WithChunkOverlap(100),
		charsplitter.WithKeepSeparator(false),
	)

	// Split text
	type WeaviateChunkObject struct {
		Href     string `json:"href"`
		Chunk    string `json:"chunk"`
		Sequence int    `json:"sequence"`
	}

	for _, hrefAndContent := range jsonMap {
		chunks, err := splitter.SplitText(hrefAndContent.Content)
		if err != nil {
			log.Printf("Failed to split text for %s, %v", hrefAndContent.Href, hrefAndContent.Content)
		}

		// Loop through each individual chunk
		for i, chunkText := range chunks {
			chunkPayload := WeaviateChunkObject{
				Href:     hrefAndContent.Href,
				Chunk:    chunkText,
				Sequence: i,
			}
			fmt.Printf("Ready for Weaviate -> Href: %s | Chunk %d/%d (Len: %d)\n",
				chunkPayload.Href,
				chunkPayload.Sequence+1,
				len(chunks),
				len(chunkPayload.Chunk),
			)
		}

	}

}
