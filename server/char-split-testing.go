package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	// Iterate
	for urlKey, values := range m {
		fmt.Printf("URL key: %s\n", urlKey)
		fmt.Printf("Stuct Href value: %s\n", values.Href)
	}

	//fmt.Println(string(data))

	//os.WriteFile("crawl_data/crawl_results_test.json", data, 0644)

	// Text splitter
	// textSplitter := charsplitter.New(
	// 	charsplitter.WithChunkSize(512),
	// 	charsplitter.WithChunkOverlap(100),
	// 	charsplitter.WithKeepSeparator(false),
	// )
}

