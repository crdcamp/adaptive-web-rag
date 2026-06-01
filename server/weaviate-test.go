package main

// Reference: https://docs.weaviate.io/weaviate/manage-collections/collection-operations

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func Test() {
	// We'll make this a struct later
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}

	//GetSchema()
	testWeaveClient, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TEST WEAVE TYPE: %T", testWeaveClient.Experimental())

	GetTestCollection(testWeaveClient, "")

	//GetTestCollection()
	//CreateTestCollection("WeavinItUp", "Let's hope this works first try", testClient)
}

func GetSchema() {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	result, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nSCHEMA: %v", result)
}

// Straight from the docs. No AI slop to be found
// Need to add check for if the collection already exists (assuming that's not done by default somehow)
func CreateTestCollection(name string, description string, client *weaviate.Client) {
	ctx := context.Background() // Still need to fully figure out how to use this
	className := name

	// This should be a function or something reusable
	emptyClass := &models.Class{
		Class:           className,
		Description:     description,
		Vectorizer:      "text2vec-openai", // Double check this. Probably wrong one
		VectorIndexType: "hnsw",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: schema.DataTypeText.PropString(),
			},
			{
				Name:     "body",
				DataType: schema.DataTypeText.PropString(),
			},
		},
	}
	// Add better error handling
	// This part of the function could also be just done a lot more elegantly in general
	// ... We'll work on that later
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(name).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if exists != true {
		err := client.Schema().ClassCreator().
			WithClass(emptyClass).
			Do(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func GetTestCollection(client *weaviate.Client, className string) {

}
