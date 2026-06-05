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
	type AllSplitText struct {
		Href      string
		SplitText string
	}

	for _, hrefAndContent := range jsonMap {
		// Create href variable
		href := hrefAndContent.Href
		// Create split text variable
		splitText, err := splitter.SplitText(hrefAndContent.Content)
		fmt.Println("splitText len:\n", len(splitText))
		if err != nil {
			panic(err)
		}

		for text := range splitText {}
		}
	}
}
