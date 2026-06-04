package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/crdcamp/charsplitter"
)

func TestSplit() {
	// Read `crawl_results.json`
	content, err := ioutil.ReadFile("crawl_data/crawl_results.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload map[string]string
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// Text splitter
	testSplitter := charsplitter.New(
		charsplitter.WithChunkSize(512),
		charsplitter.WithChunkOverlap(100),
		charsplitter.WithKeepSeparator(false),
	)

	// Format used by: https://docs.weaviate.io/weaviate/model-providers/weaviate/embeddings#configure-the-vectorizer
	// var sourceObjects = []map[string]string{
	//        {"title": "The Shawshank Redemption", "description": "A wrongfully imprisoned man forms an inspiring friendship while finding hope and redemption in the darkest of places."},
	//        {"title": "The Godfather", "description": "A powerful mafia family struggles to balance loyalty, power, and betrayal in this iconic crime saga."},
	//        {"title": "The Dark Knight", "description": "Batman faces his greatest challenge as he battles the chaos unleashed by the Joker in Gotham City."},
	//        {"title": "Jingle All the Way", "description": "A desperate father goes to hilarious lengths to secure the season's hottest toy for his son on Christmas Eve."},
	//        {"title": "A Christmas Carol", "description": "A miserly old man is transformed after being visited by three ghosts on Christmas Eve in this timeless tale of redemption."},
	//    }

	// Split a single key/value pair from `crawl_results.json`

}
