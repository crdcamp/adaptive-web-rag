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
	m := map[string]HrefContent{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		panic(err)
	}

	// Text splitter
	splitter := charsplitter.New(
		charsplitter.WithChunkSize(512),
		charsplitter.WithChunkOverlap(100),
		charsplitter.WithKeepSeparator(false),
	)

	// Iterate example
	for urlKey, hrefContent := range m {
		allSplitText := []string{}
		splitText, err := splitter.SplitText(hrefContent.Content)
		if err != nil {
			panic(err)
		}

		allSplitText.append(splitText)
	}


	// Split a single document
	}

}
