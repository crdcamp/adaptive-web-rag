package main

// Reference: https://docs.weaviate.io/weaviate/manage-collections/collection-operations

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func GetSchema() {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nSCHEMA: %v", schema)
}

func Test() {
	GetSchema()
}

func CreateTestCollection(name string, client weaviate.Client) {
	className := name

	  emptyClass := &models.Class{
	    Class: className,
	  }

	  // Create the collection (also called class)
	  err := client.Schema().ClassCreator().
	    WithClass(emptyClass).
	    Do(ctx)
