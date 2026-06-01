package main

// Reference: https://docs.weaviate.io/weaviate/manage-collections/collection-operations

import (
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

func CleanTest() {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	fmt.Printf("cdf type: %T", cfg)

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nclient type: %T", client)
}

var cfg = weaviate.Config{
	Host:   "localhost:8080",
	Scheme: "http",
}

func CreateCollection(name string, description string, client *weaviate.Client) {

}
