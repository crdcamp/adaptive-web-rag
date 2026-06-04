package main

import (
	"encoding/json"
	"log"
	"os"
)

func TestSplit() {
	// Read `crawl_results.json`
	content, err := os.ReadFile("crawl_data/crawl_results.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(m, "", " ")
	//fmt.Println(string(data))

	os.WriteFile("crawl_data/crawl_results_test.json", data, 0644)

	// Text splitter
	// textSplitter := charsplitter.New(
	// 	charsplitter.WithChunkSize(512),
	// 	charsplitter.WithChunkOverlap(100),
	// 	charsplitter.WithKeepSeparator(false),
	// )
}

