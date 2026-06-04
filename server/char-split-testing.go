package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func TestSplit() {
	// Read `crawl_results.json`
	content, err := ioutil.ReadFile("crawl_data/crawl_results.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	type CrawlData []struct {
		Href    string `json:"href"`
		Content string `json:"content"`
	}

	var c CrawlData
	err = json.Unmarshal(content, &c)
	if err != nil {
		log.Fatal("Error unmarshaling into CrawlData: ", err)
	}
	fmt.Println(c.Href)
	fmt.Println(c.Content)

	// Text splitter
	// textSplitter := charsplitter.New(
	// 	charsplitter.WithChunkSize(512),
	// 	charsplitter.WithChunkOverlap(100),
	// 	charsplitter.WithKeepSeparator(false),
	// )
}

